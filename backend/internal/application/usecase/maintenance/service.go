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
	domaincar "cvmc/internal/domain/car"
	domainmaintenance "cvmc/internal/domain/maintenance"
	domainuser "cvmc/internal/domain/user"
	"cvmc/internal/shared/filecipher"
)

var (
	ErrMaintenanceNotFound  = errors.New("maintenance not found")
	ErrMaintenanceForbidden = errors.New("forbidden")
	ErrMaintenanceInvalid   = errors.New("invalid payload")
	ErrStorageQuotaExceeded = errors.New("limite de cota de arquivos atingido (máx. 100 arquivos ou 100 MB)")
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

func (s *Service) normalizeAndEncryptAttachments(ctx context.Context, ownerID, excludeMaintID string, incoming []domainmaintenance.Attachment) ([]domainmaintenance.Attachment, error) {
	if len(incoming) == 0 {
		return incoming, nil
	}

	maxFiles := domainuser.DefaultMaxFiles
	maxStorage := domainuser.DefaultMaxStorageBytes
	if s.users != nil {
		if u, err := s.users.FindByID(ctx, ownerID); err == nil {
			maxFiles = u.EffectiveMaxFiles()
			maxStorage = u.EffectiveMaxStorageBytes()
		}
	}

	filesCount := len(incoming)
	var storageUsed int64

	prepared := make([]domainmaintenance.Attachment, len(incoming))
	for i, att := range incoming {
		realSize := filecipher.EstimateDecodedBytes(att.DataUrl)
		if realSize > att.Size {
			att.Size = realSize
		}
		if att.Size > 0 {
			storageUsed += att.Size
		}
		prepared[i] = att
	}

	if s.cars != nil {
		userCars, err := s.cars.ListByUser(ctx, ownerID)
		if err == nil {
			for _, c := range userCars {
				if c.OwnerID != ownerID {
					continue
				}
				if filecipher.IsInlineFileData(c.ImageUrl) {
					filesCount++
					storageUsed += filecipher.EstimateDecodedBytes(c.ImageUrl)
				}
				items, err := s.maintenances.ListByCar(ctx, c.ID)
				if err == nil {
					for _, m := range items {
						if m.ID == excludeMaintID {
							continue
						}
						for _, att := range m.Attachments {
							filesCount++
							sz := att.Size
							if est := filecipher.EstimateDecodedBytes(att.DataUrl); est > sz {
								sz = est
							}
							if sz > 0 {
								storageUsed += sz
							}
						}
					}
				}
			}
		}
	}

	if filesCount > maxFiles || storageUsed > maxStorage {
		return nil, ErrStorageQuotaExceeded
	}

	if s.cipher != nil {
		for i := range prepared {
			if prepared[i].DataUrl != "" {
				enc, err := s.cipher.EncryptForOwner(ownerID, prepared[i].DataUrl)
				if err != nil {
					return nil, err
				}
				prepared[i].DataUrl = enc
			}
		}
	}

	return prepared, nil
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
	preparedAttachments, err := s.normalizeAndEncryptAttachments(ctx, car.OwnerID, "", input.Attachments)
	if err != nil {
		return domainmaintenance.Maintenance{}, err
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
	if input.Mileage > car.LastMileage {
		car.LastMileage = input.Mileage
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
	preparedAttachments, err := s.normalizeAndEncryptAttachments(ctx, car.OwnerID, maintenanceID, input.Attachments)
	if err != nil {
		return domainmaintenance.Maintenance{}, err
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
	if item.Mileage > car.LastMileage {
		car.LastMileage = item.Mileage
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
	return s.maintenances.Delete(ctx, maintenanceID)
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
