package maintenance

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	carport "cvmc/internal/application/ports/car"
	maintenanceport "cvmc/internal/application/ports/maintenance"
	domainmaintenance "cvmc/internal/domain/maintenance"

	domaincar "cvmc/internal/domain/car"
)

var (
	ErrMaintenanceNotFound  = errors.New("maintenance not found")
	ErrMaintenanceForbidden = errors.New("forbidden")
	ErrMaintenanceInvalid   = errors.New("invalid payload")
)

type Service struct {
	maintenances maintenanceport.Repository
	cars         carport.Repository
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
		Attachments: input.Attachments,
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
	return created, nil
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
	sort.Slice(items, func(i, j int) bool { return items[i].Date.After(items[j].Date) })
	return items, nil
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
	item.Title = strings.TrimSpace(input.Title)
	item.Description = strings.TrimSpace(input.Description)
	if !input.Date.IsZero() {
		item.Date = input.Date
	}
	item.Mileage = input.Mileage
	item.Types = input.Types
	item.Cost = input.Cost
	item.Attachments = input.Attachments
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
	return updated, nil
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
