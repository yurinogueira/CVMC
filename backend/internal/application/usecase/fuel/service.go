package fuel

import (
	"context"
	"errors"
	"math"
	"sort"
	"strings"
	"time"

	carport "cvmc/internal/application/ports/car"
	fuelport "cvmc/internal/application/ports/fuel"
	domaincar "cvmc/internal/domain/car"
	domainfuel "cvmc/internal/domain/fuel"
)

var (
	ErrFuelingNotFound  = errors.New("fueling not found")
	ErrFuelingForbidden = errors.New("forbidden")
	ErrFuelingInvalid   = errors.New("invalid payload")
)

type Service struct {
	fuelings fuelport.Repository
	cars     carport.Repository
	now      func() time.Time
}

type CreateInput struct {
	Date          time.Time
	FuelType      string
	Liters        float64
	PricePerLiter float64
	TotalCost     float64
	IsFullTank    bool
	GasStation    string
	Notes         string
}

type UpdateInput struct {
	Date          time.Time
	FuelType      string
	Liters        float64
	PricePerLiter float64
	TotalCost     float64
	IsFullTank    bool
	GasStation    string
	Notes         string
}

func NewService(fuelings fuelport.Repository, cars carport.Repository) *Service {
	return &Service{fuelings: fuelings, cars: cars, now: time.Now}
}

func (s *Service) Create(ctx context.Context, actorID, carID string, input CreateInput) (domainfuel.Fueling, error) {
	car, err := s.cars.GetByID(ctx, carID)
	if err != nil {
		return domainfuel.Fueling{}, ErrFuelingNotFound
	}
	if car.OwnerID != actorID {
		return domainfuel.Fueling{}, ErrFuelingForbidden
	}

	fuelType := strings.TrimSpace(input.FuelType)
	if fuelType == "" || len(fuelType) > 50 {
		return domainfuel.Fueling{}, ErrFuelingInvalid
	}
	if input.Liters <= 0 {
		return domainfuel.Fueling{}, ErrFuelingInvalid
	}
	if len(input.GasStation) > 100 || len(input.Notes) > 500 {
		return domainfuel.Fueling{}, ErrFuelingInvalid
	}

	pricePerLiter := input.PricePerLiter
	totalCost := input.TotalCost

	if totalCost <= 0 && pricePerLiter > 0 {
		totalCost = math.Round(input.Liters*pricePerLiter*100) / 100
	} else if pricePerLiter <= 0 && totalCost > 0 {
		pricePerLiter = math.Round((totalCost/input.Liters)*1000) / 1000
	}

	if totalCost <= 0 || pricePerLiter <= 0 {
		return domainfuel.Fueling{}, ErrFuelingInvalid
	}

	when := input.Date
	if when.IsZero() {
		when = s.now().UTC()
	}

	item := domainfuel.Fueling{
		CarID:         carID,
		Date:          when,
		FuelType:      fuelType,
		Liters:        input.Liters,
		PricePerLiter: pricePerLiter,
		TotalCost:     totalCost,
		IsFullTank:    input.IsFullTank,
		GasStation:    strings.TrimSpace(input.GasStation),
		Notes:         strings.TrimSpace(input.Notes),
		CreatedAt:     s.now().UTC(),
		UpdatedAt:     s.now().UTC(),
	}

	created, err := s.fuelings.Create(ctx, item)
	if err != nil {
		return domainfuel.Fueling{}, err
	}

	return created, nil
}

func (s *Service) List(ctx context.Context, actorID, carID string) ([]domainfuel.Fueling, error) {
	car, err := s.cars.GetByID(ctx, carID)
	if err != nil {
		return nil, ErrFuelingNotFound
	}
	if !accessible(actorID, car) {
		return nil, ErrFuelingForbidden
	}

	items, err := s.fuelings.ListByCar(ctx, carID)
	if err != nil {
		return nil, err
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Date.After(items[j].Date)
	})

	return items, nil
}

func (s *Service) Get(ctx context.Context, actorID, fuelingID string) (domainfuel.Fueling, error) {
	item, err := s.fuelings.GetByID(ctx, fuelingID)
	if err != nil {
		return domainfuel.Fueling{}, ErrFuelingNotFound
	}
	car, err := s.cars.GetByID(ctx, item.CarID)
	if err != nil {
		return domainfuel.Fueling{}, ErrFuelingNotFound
	}
	if !accessible(actorID, car) {
		return domainfuel.Fueling{}, ErrFuelingForbidden
	}
	return item, nil
}

func (s *Service) Update(ctx context.Context, actorID, fuelingID string, input UpdateInput) (domainfuel.Fueling, error) {
	item, err := s.fuelings.GetByID(ctx, fuelingID)
	if err != nil {
		return domainfuel.Fueling{}, ErrFuelingNotFound
	}
	car, err := s.cars.GetByID(ctx, item.CarID)
	if err != nil {
		return domainfuel.Fueling{}, ErrFuelingNotFound
	}
	if car.OwnerID != actorID {
		return domainfuel.Fueling{}, ErrFuelingForbidden
	}

	fuelType := strings.TrimSpace(input.FuelType)
	if fuelType == "" || len(fuelType) > 50 {
		return domainfuel.Fueling{}, ErrFuelingInvalid
	}
	if input.Liters <= 0 {
		return domainfuel.Fueling{}, ErrFuelingInvalid
	}
	if len(input.GasStation) > 100 || len(input.Notes) > 500 {
		return domainfuel.Fueling{}, ErrFuelingInvalid
	}

	pricePerLiter := input.PricePerLiter
	totalCost := input.TotalCost

	if totalCost <= 0 && pricePerLiter > 0 {
		totalCost = math.Round(input.Liters*pricePerLiter*100) / 100
	} else if pricePerLiter <= 0 && totalCost > 0 {
		pricePerLiter = math.Round((totalCost/input.Liters)*1000) / 1000
	}

	if totalCost <= 0 || pricePerLiter <= 0 {
		return domainfuel.Fueling{}, ErrFuelingInvalid
	}

	if !input.Date.IsZero() {
		item.Date = input.Date
	}
	item.FuelType = fuelType
	item.Liters = input.Liters
	item.PricePerLiter = pricePerLiter
	item.TotalCost = totalCost
	item.IsFullTank = input.IsFullTank
	item.GasStation = strings.TrimSpace(input.GasStation)
	item.Notes = strings.TrimSpace(input.Notes)
	item.UpdatedAt = s.now().UTC()

	updated, err := s.fuelings.Update(ctx, item)
	if err != nil {
		return domainfuel.Fueling{}, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, actorID, fuelingID string) error {
	item, err := s.fuelings.GetByID(ctx, fuelingID)
	if err != nil {
		return ErrFuelingNotFound
	}
	car, err := s.cars.GetByID(ctx, item.CarID)
	if err != nil {
		return ErrFuelingNotFound
	}
	if car.OwnerID != actorID {
		return ErrFuelingForbidden
	}
	return s.fuelings.Delete(ctx, fuelingID)
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
