package maintenance

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
	domainmaintenance "cvmc/internal/domain/maintenance"
	domainuser "cvmc/internal/domain/user"
	"cvmc/internal/shared/filecipher"
)

var (
	ErrMaintenanceNotFound  = errors.New("maintenance not found")
	ErrMaintenanceForbidden = errors.New("forbidden")
	ErrMaintenanceInvalid   = errors.New("invalid payload")
	ErrStorageQuotaExceeded = storagequota.ErrStorageQuotaExceeded
)

type Service struct {
	maintenances maintenanceport.Repository
	cars         carport.Repository
	users        userport.Repository
	cipher       *filecipher.FileCipher
	now          func() time.Time
}

type CreateInput struct {
	Title       string
	Description string
	Date        time.Time
	Mileage     int
	Types       []string
	Cost        *float64
	Attachments []domainmaintenance.Attachment
}

type UpdateInput struct {
	Title       string
	Description string
	Date        time.Time
	Mileage     int
	Types       []string
	Cost        *float64
	Attachments []domainmaintenance.Attachment
}

func NewService(maintenances maintenanceport.Repository, cars carport.Repository) *Service {
	return &Service{maintenances: maintenances, cars: cars, now: time.Now}
}

func (s *Service) SetSecurityDependencies(users userport.Repository, cipher *filecipher.FileCipher) {
	s.users = users
	s.cipher = cipher
}

func (s *Service) prepareAndEncryptAttachments(ownerID string, incoming []domainmaintenance.Attachment) ([]domainmaintenance.Attachment, int, int64, error) {
	if len(incoming) == 0 {
		return incoming, 0, 0, nil
	}

	var totalBytes int64
	prepared := make([]domainmaintenance.Attachment, len(incoming))
	for i, att := range incoming {
		realSize := filecipher.EstimateDecodedBytes(att.DataUrl)
		if realSize > att.Size {
			att.Size = realSize
		}
		if att.Size > 0 {
			totalBytes += att.Size
		}
		prepared[i] = att
	}

	if s.cipher != nil {
		for i := range prepared {
			if prepared[i].DataUrl != "" {
				enc, err := s.cipher.EncryptForOwner(ownerID, prepared[i].DataUrl)
				if err != nil {
					return nil, 0, 0, err
				}
				prepared[i].DataUrl = enc
			}
		}
	}

	return prepared, len(prepared), totalBytes, nil
}

func (s *Service) decryptMaintenance(ownerID string, m domainmaintenance.Maintenance) domainmaintenance.Maintenance {
	if s.cipher == nil || len(m.Attachments) == 0 {
		return m
	}
	decrypted := make([]domainmaintenance.Attachment, len(m.Attachments))
	copy(decrypted, m.Attachments)
	for i := range decrypted {
		if decrypted[i].DataUrl != "" {
			if plain, err := s.cipher.DecryptForOwner(ownerID, decrypted[i].DataUrl); err == nil {
				decrypted[i].DataUrl = plain
			}
		}
	}
	m.Attachments = decrypted
	return m
}

func (s *Service) Create(ctx context.Context, actorID, carID string, input CreateInput) (domainmaintenance.Maintenance, error) {
	car, err := s.cars.GetByID(ctx, carID)
	if err != nil {
		return domainmaintenance.Maintenance{}, ErrMaintenanceNotFound
	}
	if car.OwnerID != actorID {
		return domainmaintenance.Maintenance{}, ErrMaintenanceForbidden
	}
	if strings.TrimSpace(input.Title) == "" {
		return domainmaintenance.Maintenance{}, ErrMaintenanceInvalid
	}
	preparedAttachments, newFiles, newBytes, err := s.prepareAndEncryptAttachments(car.OwnerID, input.Attachments)
	if err != nil {
		return domainmaintenance.Maintenance{}, err
	}

	var syncedUser domainuser.User
	var hasSyncedUser bool
	if s.users != nil {
		if u, err := s.users.FindByID(ctx, car.OwnerID); err == nil {
			u, _ = storagequota.EnsureSynced(ctx, u, s.users, s.cars, s.maintenances, s.cipher)
			if err := storagequota.CheckDelta(u, newFiles, newBytes); err != nil {
				return domainmaintenance.Maintenance{}, ErrStorageQuotaExceeded
			}
			syncedUser = u
			hasSyncedUser = true
		}
	} else if newFiles > domainuser.DefaultMaxFiles || newBytes > domainuser.DefaultMaxStorageBytes {
		return domainmaintenance.Maintenance{}, ErrStorageQuotaExceeded
	}

	when := input.Date
	if when.IsZero() {
		when = s.now().UTC()
	}
	maintenance := domainmaintenance.Maintenance{
		CarID:       carID,
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Date:        when,
		Mileage:     input.Mileage,
		Types:       input.Types,
		Cost:        input.Cost,
		Attachments: preparedAttachments,
		CreatedAt:   s.now().UTC(),
		UpdatedAt:   s.now().UTC(),
	}
	created, err := s.maintenances.Create(ctx, maintenance)
	if err != nil {
		return domainmaintenance.Maintenance{}, err
	}
	if hasSyncedUser && (newFiles != 0 || newBytes != 0) {
		_, _ = storagequota.CommitDelta(ctx, syncedUser, s.users, newFiles, newBytes)
	}
	if input.Mileage > car.LastMileage {
		car.LastMileage = input.Mileage
		car.MileageUpdatedAt = s.now().UTC()
		car.UpdatedAt = s.now().UTC()
		_, _ = s.cars.Update(ctx, car)
	}
	return s.decryptMaintenance(car.OwnerID, created), nil
}

func (s *Service) List(ctx context.Context, actorID, carID string) ([]domainmaintenance.Maintenance, error) {
	car, err := s.cars.GetByID(ctx, carID)
	if err != nil {
		return nil, ErrMaintenanceNotFound
	}
	if !accessible(actorID, car) {
		return nil, ErrMaintenanceForbidden
	}
	items, err := s.maintenances.ListByCar(ctx, carID)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i] = s.decryptMaintenance(car.OwnerID, items[i])
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Date.After(items[j].Date) })
	return items, nil
}

func (s *Service) Get(ctx context.Context, actorID, maintenanceID string) (domainmaintenance.Maintenance, error) {
	item, err := s.maintenances.GetByID(ctx, maintenanceID)
	if err != nil {
		return domainmaintenance.Maintenance{}, ErrMaintenanceNotFound
	}
	car, err := s.cars.GetByID(ctx, item.CarID)
	if err != nil {
		return domainmaintenance.Maintenance{}, ErrMaintenanceNotFound
	}
	if !accessible(actorID, car) {
		return domainmaintenance.Maintenance{}, ErrMaintenanceForbidden
	}
	return s.decryptMaintenance(car.OwnerID, item), nil
}

func (s *Service) Update(ctx context.Context, actorID, maintenanceID string, input UpdateInput) (domainmaintenance.Maintenance, error) {
	item, err := s.maintenances.GetByID(ctx, maintenanceID)
	if err != nil {
		return domainmaintenance.Maintenance{}, ErrMaintenanceNotFound
	}
	car, err := s.cars.GetByID(ctx, item.CarID)
	if err != nil {
		return domainmaintenance.Maintenance{}, ErrMaintenanceNotFound
	}
	if car.OwnerID != actorID {
		return domainmaintenance.Maintenance{}, ErrMaintenanceForbidden
	}
	if strings.TrimSpace(input.Title) == "" {
		return domainmaintenance.Maintenance{}, ErrMaintenanceInvalid
	}

	oldFiles, oldBytes := storagequota.AttachmentMetrics(item.Attachments)
	preparedAttachments, newFiles, newBytes, err := s.prepareAndEncryptAttachments(car.OwnerID, input.Attachments)
	if err != nil {
		return domainmaintenance.Maintenance{}, err
	}
	deltaFiles := newFiles - oldFiles
	deltaBytes := newBytes - oldBytes

	var syncedUser domainuser.User
	var hasSyncedUser bool
	if s.users != nil {
		if u, err := s.users.FindByID(ctx, car.OwnerID); err == nil {
			u, _ = storagequota.EnsureSynced(ctx, u, s.users, s.cars, s.maintenances, s.cipher)
			if err := storagequota.CheckDelta(u, deltaFiles, deltaBytes); err != nil {
				return domainmaintenance.Maintenance{}, ErrStorageQuotaExceeded
			}
			syncedUser = u
			hasSyncedUser = true
		}
	} else if newFiles > domainuser.DefaultMaxFiles || newBytes > domainuser.DefaultMaxStorageBytes {
		return domainmaintenance.Maintenance{}, ErrStorageQuotaExceeded
	}

	item.Title = strings.TrimSpace(input.Title)
	item.Description = strings.TrimSpace(input.Description)
	if !input.Date.IsZero() {
		item.Date = input.Date
	}
	item.Mileage = input.Mileage
	item.Types = input.Types
	item.Cost = input.Cost
	item.Attachments = preparedAttachments
	item.UpdatedAt = s.now().UTC()
	updated, err := s.maintenances.Update(ctx, item)
	if err != nil {
		return domainmaintenance.Maintenance{}, err
	}
	if hasSyncedUser && (deltaFiles != 0 || deltaBytes != 0) {
		_, _ = storagequota.CommitDelta(ctx, syncedUser, s.users, deltaFiles, deltaBytes)
	}
	if item.Mileage > car.LastMileage {
		car.LastMileage = item.Mileage
		car.MileageUpdatedAt = s.now().UTC()
		car.UpdatedAt = s.now().UTC()
		_, _ = s.cars.Update(ctx, car)
	}
	return s.decryptMaintenance(car.OwnerID, updated), nil
}

func (s *Service) Delete(ctx context.Context, actorID, maintenanceID string) error {
	item, err := s.maintenances.GetByID(ctx, maintenanceID)
	if err != nil {
		return ErrMaintenanceNotFound
	}
	car, err := s.cars.GetByID(ctx, item.CarID)
	if err != nil {
		return ErrMaintenanceNotFound
	}
	if car.OwnerID != actorID {
		return ErrMaintenanceForbidden
	}

	var syncedUser domainuser.User
	var hasSyncedUser bool
	if s.users != nil {
		if u, err := s.users.FindByID(ctx, car.OwnerID); err == nil {
			u, _ = storagequota.EnsureSynced(ctx, u, s.users, s.cars, s.maintenances, s.cipher)
			syncedUser = u
			hasSyncedUser = true
		}
	}
	oldFiles, oldBytes := storagequota.AttachmentMetrics(item.Attachments)

	if err := s.maintenances.Delete(ctx, maintenanceID); err != nil {
		return err
	}
	if hasSyncedUser && (oldFiles != 0 || oldBytes != 0) {
		_, _ = storagequota.CommitDelta(ctx, syncedUser, s.users, -oldFiles, -oldBytes)
	}
	return nil
}

func accessible(actorID string, car domaincar.Car) bool {
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
