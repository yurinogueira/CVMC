package car

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	carport "cvmc/internal/application/ports/car"
	userport "cvmc/internal/application/ports/user"
	domaincar "cvmc/internal/domain/car"
	memoryuser "cvmc/internal/infrastructure/user/memory"
)

var (
	ErrCarNotFound    = errors.New("car not found")
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidPayload = errors.New("invalid payload")
	ErrUserNotFound   = errors.New("user not found")
	ErrShareNotFound  = errors.New("share target not found")
)

type Service struct {
	cars  carport.Repository
	users userport.Repository
	now   func() time.Time
}

type CreateInput struct {
	Name            string
	Manufacturer    string
	Model           string
	YearManufacture int
	YearModel       int
	LastMileage     int
}

type UpdateInput struct {
	Name            string
	Manufacturer    string
	Model           string
	YearManufacture int
	YearModel       int
	LastMileage     int
}

func NewService(cars carport.Repository, users userport.Repository) *Service {
	return &Service{cars: cars, users: users, now: time.Now}
}

func (s *Service) Create(ctx context.Context, ownerID string, input CreateInput) (domaincar.Car, error) {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Manufacturer) == "" || strings.TrimSpace(input.Model) == "" {
		return domaincar.Car{}, ErrInvalidPayload
	}
	car := domaincar.Car{
		OwnerID:         ownerID,
		Name:            strings.TrimSpace(input.Name),
		Manufacturer:    strings.TrimSpace(input.Manufacturer),
		Model:           strings.TrimSpace(input.Model),
		YearManufacture: input.YearManufacture,
		YearModel:       input.YearModel,
		LastMileage:     input.LastMileage,
		SharedWith:      []string{},
		CreatedAt:       s.now().UTC(),
		UpdatedAt:       s.now().UTC(),
	}
	return s.cars.Create(ctx, car)
}

func (s *Service) Get(ctx context.Context, actorID, carID string) (domaincar.Car, error) {
	car, err := s.cars.GetByID(ctx, carID)
	if err != nil {
		return domaincar.Car{}, ErrCarNotFound
	}
	if !canAccessCar(actorID, car) {
		return domaincar.Car{}, ErrForbidden
	}
	return car, nil
}

func (s *Service) List(ctx context.Context, actorID string) ([]domaincar.Car, error) {
	items, err := s.cars.ListByUser(ctx, actorID)
	if err != nil {
		return nil, err
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
	car.Name = strings.TrimSpace(input.Name)
	car.Manufacturer = strings.TrimSpace(input.Manufacturer)
	car.Model = strings.TrimSpace(input.Model)
	car.YearManufacture = input.YearManufacture
	car.YearModel = input.YearModel
	car.LastMileage = input.LastMileage
	car.UpdatedAt = s.now().UTC()
	return s.cars.Update(ctx, car)
}

func (s *Service) Delete(ctx context.Context, actorID, carID string) error {
	car, err := s.cars.GetByID(ctx, carID)
	if err != nil {
		return ErrCarNotFound
	}
	if car.OwnerID != actorID {
		return ErrForbidden
	}
	return s.cars.Delete(ctx, carID)
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
		if errors.Is(err, memoryuser.ErrNotFound) {
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
