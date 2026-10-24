package maintenance

import (
	"context"
	"sort"

	domainmaintenance "cvmc/internal/domain/maintenance"
)

// GetUpcomingByCar calculates the upcoming periodic maintenances for a specific vehicle.
func (s *Service) GetUpcomingByCar(ctx context.Context, actorID, carID string) (domainmaintenance.CarUpcomingMaintenances, error) {
	car, err := s.cars.GetByID(ctx, carID)
	if err != nil {
		return domainmaintenance.CarUpcomingMaintenances{}, ErrMaintenanceNotFound
	}
	if !accessible(actorID, car) {
		return domainmaintenance.CarUpcomingMaintenances{}, ErrMaintenanceForbidden
	}

	items, err := s.maintenances.ListByCar(ctx, carID)
	if err != nil {
		return domainmaintenance.CarUpcomingMaintenances{}, err
	}

	return domainmaintenance.CalculateCarUpcoming(car, items, s.now().UTC()), nil
}

// GetUpcomingSummary calculates the consolidated upcoming maintenances across all vehicles accessible by the user.
func (s *Service) GetUpcomingSummary(ctx context.Context, actorID string) (domainmaintenance.FleetUpcomingSummary, error) {
	cars, err := s.cars.ListByUser(ctx, actorID)
	if err != nil {
		return domainmaintenance.FleetUpcomingSummary{}, err
	}

	sort.Slice(cars, func(i, j int) bool {
		return cars[i].CreatedAt.After(cars[j].CreatedAt)
	})

	now := s.now().UTC()
	vehicles := make([]domainmaintenance.CarUpcomingMaintenances, 0, len(cars))
	totalOK := 0
	totalUpcoming := 0
	totalOverdue := 0

	for _, car := range cars {
		items, err := s.maintenances.ListByCar(ctx, car.ID)
		if err != nil {
			items = nil
		}
		evaluated := domainmaintenance.CalculateCarUpcoming(car, items, now)
		totalOK += evaluated.OKCount
		totalUpcoming += evaluated.UpcomingCount
		totalOverdue += evaluated.OverdueCount
		vehicles = append(vehicles, evaluated)
	}

	fleetStatus := domainmaintenance.StatusOK
	if totalOverdue > 0 {
		fleetStatus = domainmaintenance.StatusOverdue
	} else if totalUpcoming > 0 {
		fleetStatus = domainmaintenance.StatusUpcoming
	}

	return domainmaintenance.FleetUpcomingSummary{
		TotalVehicles: len(cars),
		Status:        fleetStatus,
		OKCount:       totalOK,
		UpcomingCount: totalUpcoming,
		OverdueCount:  totalOverdue,
		Vehicles:      vehicles,
	}, nil
}
