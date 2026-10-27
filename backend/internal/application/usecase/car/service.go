package car

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	carport "cvmc/internal/application/ports/car"
	maintenanceport "cvmc/internal/application/ports/maintenance"
	userport "cvmc/internal/application/ports/user"
	"cvmc/internal/application/usecase/storagequota"
	domaincar "cvmc/internal/domain/car"
	domainuser "cvmc/internal/domain/user"
	memoryuser "cvmc/internal/infrastructure/user/memory"
	"cvmc/internal/shared/filecipher"
)

var (
	ErrCarNotFound          = errors.New("car not found")
	ErrForbidden            = errors.New("forbidden")
	ErrInvalidPayload       = errors.New("invalid payload")
	ErrUserNotFound         = errors.New("user not found")
	ErrShareNotFound        = errors.New("share target not found")
	ErrEmailNotVerified     = errors.New("valide seu e-mail para cadastrar veículos")
	ErrVehicleLimitReached  = errors.New("limite de veículos atingido")
	ErrStorageQuotaExceeded = storagequota.ErrStorageQuotaExceeded
)

type Service struct {
	cars         carport.Repository
	users        userport.Repository
	maintenances maintenanceport.Repository
	cipher       *filecipher.FileCipher
	now          func() time.Time
}

type CreateInput struct {
	Name            string
	Manufacturer    string
	Model           string
	YearManufacture int
	YearModel       int
	LastMileage     int
	VehicleType     string
	ImageUrl        string
	FIPECode        string
	FIPEPrice       string
	Fuel            string
}

type UpdateInput struct {
	Name                 string
	Manufacturer         string
	Model                string
	YearManufacture      int
	YearModel            int
	LastMileage          int
	VehicleType          string
	ImageUrl             string
	FIPECode             string
	FIPEPrice            string
	Fuel                 string
	MaintenanceIntervals map[string]domaincar.MaintenanceIntervalConfig
}

func NewService(cars carport.Repository, users userport.Repository) *Service {
	return &Service{cars: cars, users: users, now: time.Now}
}

func (s *Service) SetSecurityDependencies(maintenances maintenanceport.Repository, cipher *filecipher.FileCipher) {
	s.maintenances = maintenances
	s.cipher = cipher
}

func (s *Service) decryptCar(car domaincar.Car) domaincar.Car {
	if s.cipher != nil && car.ImageUrl != "" {
		if plain, err := s.cipher.DecryptForOwner(car.OwnerID, car.ImageUrl); err == nil {
			car.ImageUrl = plain
		}
	}
	return car
}

func (s *Service) Create(ctx context.Context, ownerID string, input CreateInput) (domaincar.Car, error) {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Manufacturer) == "" || strings.TrimSpace(input.Model) == "" {
		return domaincar.Car{}, ErrInvalidPayload
	}

	user, err := s.users.FindByID(ctx, ownerID)
	if err != nil {
		return domaincar.Car{}, ErrUserNotFound
	}

	if !user.EmailVerified {
		return domaincar.Car{}, ErrEmailNotVerified
	}

	user, _ = storagequota.EnsureSynced(ctx, user, s.users, s.cars, s.maintenances, s.cipher)

	userCars, err := s.cars.ListByUser(ctx, ownerID)
	if err != nil {
		return domaincar.Car{}, err
	}

	ownedCount := 0
	for _, c := range userCars {
		if c.OwnerID == ownerID {
			ownedCount++
		}
	}

	maxVehicles := user.MaxVehicles
	if maxVehicles <= 0 {
		maxVehicles = 3
	}

	if ownedCount >= maxVehicles {
		return domaincar.Car{}, ErrVehicleLimitReached
	}

	trimmedImg := strings.TrimSpace(input.ImageUrl)
	deltaFiles, deltaBytes := storagequota.ImageMetrics(trimmedImg)
	if err := storagequota.CheckDelta(user, deltaFiles, deltaBytes); err != nil {
		return domaincar.Car{}, ErrStorageQuotaExceeded
	}
	storedImg := trimmedImg
	if s.cipher != nil && storedImg != "" {
		encImg, err := s.cipher.EncryptForOwner(ownerID, storedImg)
		if err != nil {
			return domaincar.Car{}, err
		}
		storedImg = encImg
	}

	car := domaincar.Car{
		OwnerID:         ownerID,
		Name:            strings.TrimSpace(input.Name),
		Manufacturer:    strings.TrimSpace(input.Manufacturer),
		Model:           strings.TrimSpace(input.Model),
		YearManufacture: input.YearManufacture,
		YearModel:       input.YearModel,
		InitialMileage:  input.LastMileage,
		LastMileage:     input.LastMileage,
		VehicleType:     strings.TrimSpace(input.VehicleType),
		ImageUrl:        storedImg,
		FIPECode:        strings.TrimSpace(input.FIPECode),
		FIPEPrice:       strings.TrimSpace(input.FIPEPrice),
		Fuel:            strings.TrimSpace(input.Fuel),
		SharedWith:      []string{},
		CreatedAt:       s.now().UTC(),
		UpdatedAt:       s.now().UTC(),
	}
	created, err := s.cars.Create(ctx, car)
	if err != nil {
		return domaincar.Car{}, err
	}
	if deltaFiles != 0 || deltaBytes != 0 {
		_, _ = storagequota.CommitDelta(ctx, user, s.users, deltaFiles, deltaBytes)
	}
	return s.decryptCar(created), nil
}

func (s *Service) Get(ctx context.Context, actorID, carID string) (domaincar.Car, error) {
	car, err := s.cars.GetByID(ctx, carID)
	if err != nil {
		return domaincar.Car{}, ErrCarNotFound
	}
	if !canAccessCar(actorID, car) {
		return domaincar.Car{}, ErrForbidden
	}
	return s.decryptCar(car), nil
}

func (s *Service) List(ctx context.Context, actorID string) ([]domaincar.Car, error) {
	if s.users != nil {
		if user, err := s.users.FindByID(ctx, actorID); err == nil && !user.StorageQuotaSynced {
			_, _ = storagequota.EnsureSynced(ctx, user, s.users, s.cars, s.maintenances, s.cipher)
		}
	}
	items, err := s.cars.ListByUser(ctx, actorID)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i] = s.decryptCar(items[i])
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	return items, nil
}

func (s *Service) Update(ctx context.Context, actorID, carID string, input UpdateInput) (domaincar.Car, error) {
	car, err := s.cars.GetByID(ctx, carID)
	if err != nil {
		return domaincar.Car{}, ErrCarNotFound
	}
	if car.OwnerID != actorID {
		return domaincar.Car{}, ErrForbidden
	}
	if input.LastMileage < 0 {
		return domaincar.Car{}, ErrInvalidPayload
	}
	if strings.TrimSpace(input.Name) != "" {
		car.Name = strings.TrimSpace(input.Name)
	}
	if strings.TrimSpace(input.Manufacturer) != "" {
		car.Manufacturer = strings.TrimSpace(input.Manufacturer)
	}
	if strings.TrimSpace(input.Model) != "" {
		car.Model = strings.TrimSpace(input.Model)
	}
	if input.YearManufacture > 0 {
		car.YearManufacture = input.YearManufacture
	}
	if input.YearModel > 0 {
		car.YearModel = input.YearModel
	}
	if input.LastMileage >= 0 {
		car.LastMileage = input.LastMileage
	}
	if input.VehicleType != "" {
		car.VehicleType = strings.TrimSpace(input.VehicleType)
	}

	var deltaFiles int
	var deltaBytes int64
	var syncedUser domainuser.User
	var hasSyncedUser bool

	if input.ImageUrl != "" {
		trimmedImg := strings.TrimSpace(input.ImageUrl)
		oldFiles, oldBytes := storagequota.ImageMetrics(car.ImageUrl)
		newFiles, newBytes := storagequota.ImageMetrics(trimmedImg)
		deltaFiles = newFiles - oldFiles
		deltaBytes = newBytes - oldBytes

		if s.users != nil {
			if u, err := s.users.FindByID(ctx, actorID); err == nil {
				u, _ = storagequota.EnsureSynced(ctx, u, s.users, s.cars, s.maintenances, s.cipher)
				if err := storagequota.CheckDelta(u, deltaFiles, deltaBytes); err != nil {
					return domaincar.Car{}, ErrStorageQuotaExceeded
				}
				syncedUser = u
				hasSyncedUser = true
			}
		}
		if s.cipher != nil {
			encImg, err := s.cipher.EncryptForOwner(car.OwnerID, trimmedImg)
			if err != nil {
				return domaincar.Car{}, err
			}
			car.ImageUrl = encImg
		} else {
			car.ImageUrl = trimmedImg
		}
	}
	if input.FIPECode != "" {
		car.FIPECode = strings.TrimSpace(input.FIPECode)
	}
	if input.FIPEPrice != "" {
		car.FIPEPrice = strings.TrimSpace(input.FIPEPrice)
	}
	if input.Fuel != "" {
		car.Fuel = strings.TrimSpace(input.Fuel)
	}
	if input.MaintenanceIntervals != nil {
		car.MaintenanceIntervals = input.MaintenanceIntervals
	}
	car.UpdatedAt = s.now().UTC()
	updated, err := s.cars.Update(ctx, car)
	if err != nil {
		return domaincar.Car{}, err
	}
	if hasSyncedUser && (deltaFiles != 0 || deltaBytes != 0) {
		_, _ = storagequota.CommitDelta(ctx, syncedUser, s.users, deltaFiles, deltaBytes)
	}
	return s.decryptCar(updated), nil
}

func (s *Service) Delete(ctx context.Context, actorID, carID string) error {
	car, err := s.cars.GetByID(ctx, carID)
	if err != nil {
		return ErrCarNotFound
	}
	if car.OwnerID != actorID {
		return ErrForbidden
	}

	var syncedUser domainuser.User
	var hasSyncedUser bool
	if s.users != nil {
		if u, err := s.users.FindByID(ctx, actorID); err == nil {
			u, _ = storagequota.EnsureSynced(ctx, u, s.users, s.cars, s.maintenances, s.cipher)
			syncedUser = u
			hasSyncedUser = true
		}
	}

	removedFiles, removedBytes := storagequota.ImageMetrics(car.ImageUrl)
	if s.maintenances != nil {
		if items, err := s.maintenances.ListByCar(ctx, carID); err == nil {
			for _, m := range items {
				mf, mb := storagequota.AttachmentMetrics(m.Attachments)
				removedFiles += mf
				removedBytes += mb
			}
		}
	}

	if err := s.cars.Delete(ctx, carID); err != nil {
		return err
	}
	if hasSyncedUser && (removedFiles != 0 || removedBytes != 0) {
		_, _ = storagequota.CommitDelta(ctx, syncedUser, s.users, -removedFiles, -removedBytes)
	}
	return nil
}

func (s *Service) Share(ctx context.Context, actorID, carID, email string) (domaincar.Car, error) {
	car, err := s.cars.GetByID(ctx, carID)
	if err != nil {
		return domaincar.Car{}, ErrCarNotFound
	}
	if car.OwnerID != actorID {
		return domaincar.Car{}, ErrForbidden
	}
	target, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, userport.ErrNotFound) || errors.Is(err, memoryuser.ErrNotFound) {
			return domaincar.Car{}, ErrShareNotFound
		}
		return domaincar.Car{}, err
	}
	if target.ID == car.OwnerID {
		return domaincar.Car{}, ErrInvalidPayload
	}
	updated, err := s.cars.Share(ctx, carID, target.ID)
	if err != nil {
		return domaincar.Car{}, err
	}
	return updated, nil
}

func (s *Service) Unshare(ctx context.Context, actorID, carID, userID string) (domaincar.Car, error) {
	car, err := s.cars.GetByID(ctx, carID)
	if err != nil {
		return domaincar.Car{}, ErrCarNotFound
	}
	if car.OwnerID != actorID {
		return domaincar.Car{}, ErrForbidden
	}
	return s.cars.Unshare(ctx, carID, userID)
}

func canAccessCar(actorID string, car domaincar.Car) bool {
	if car.OwnerID == actorID {
		return true
	}
	for _, shared := range car.SharedWith {
		if shared == actorID {
			return true
		}
	}
	return false
}
