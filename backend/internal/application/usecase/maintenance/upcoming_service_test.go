package maintenance

import (
	"context"
	"errors"
	"testing"
	"time"

	domaincar "cvmc/internal/domain/car"
	domainmaintenance "cvmc/internal/domain/maintenance"
	carMemory "cvmc/internal/infrastructure/car/memory"
	maintMemory "cvmc/internal/infrastructure/maintenance/memory"
)

func TestService_GetUpcomingByCar_CarAndMotorcycleIntervals(t *testing.T) {
	ctx := context.Background()
	carRepo := carMemory.NewRepository()
	maintRepo := maintMemory.NewRepository()
	svc := NewService(maintRepo, carRepo)

	fixedNow := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return fixedNow }

	// 1. Create a Car ("cars") with initial 10,000 km and current 19,200 km
	sedan, err := carRepo.Create(ctx, domaincar.Car{
		OwnerID:        "user-1",
		Name:           "Toyota Corolla",
		Manufacturer:   "Toyota",
		Model:          "Corolla XEi",
		InitialMileage: 10000,
		LastMileage:    19200,
		VehicleType:    "cars",
		CreatedAt:      fixedNow.AddDate(0, -6, 0),
	})
	if err != nil {
		t.Fatalf("failed to create car: %v", err)
	}

	// Register oil change at 10,000 km (6 months ago) -> next due at 20,000 km -> 800 km remaining (Upcoming!)
	_, err = maintRepo.Create(ctx, domainmaintenance.Maintenance{
		CarID:   sedan.ID,
		Title:   "Troca de óleo",
		Date:    fixedNow.AddDate(0, -6, 0),
		Mileage: 10000,
		Types:   []string{domainmaintenance.TypeEngineOil},
	})
	if err != nil {
		t.Fatalf("failed to create maintenance: %v", err)
	}

	upcomingSedan, err := svc.GetUpcomingByCar(ctx, "user-1", sedan.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(upcomingSedan.Items) != 5 {
		t.Fatalf("expected 5 periodic items for car, got %d", len(upcomingSedan.Items))
	}

	// Engine oil should be "upcoming" (800 km remaining <= 1000 km)
	oilItem := upcomingSedan.Items[0]
	if oilItem.Type != domainmaintenance.TypeEngineOil {
		t.Errorf("expected first item to be %q, got %q", domainmaintenance.TypeEngineOil, oilItem.Type)
	}
	if !oilItem.HasHistory {
		t.Errorf("expected engine oil to have history")
	}
	if oilItem.RemainingKm != 800 {
		t.Errorf("expected 800 km remaining, got %d", oilItem.RemainingKm)
	}
	if oilItem.Status != domainmaintenance.StatusUpcoming {
		t.Errorf("expected status upcoming, got %q", oilItem.Status)
	}

	// Cabin filter (20,000 km interval from baseline 10,000 -> next due 30,000 -> 10,800 km remaining -> OK)
	var cabinItem domainmaintenance.UpcomingMaintenanceItem
	for _, it := range upcomingSedan.Items {
		if it.Type == domainmaintenance.TypeCabinFilter {
			cabinItem = it
		}
	}
	if cabinItem.IntervalKm != 20000 {
		t.Errorf("expected 20000 km interval for cabin filter, got %d", cabinItem.IntervalKm)
	}
	if cabinItem.Status != domainmaintenance.StatusOK {
		t.Errorf("expected cabin filter status ok, got %q", cabinItem.Status)
	}

	// 2. Create a Motorcycle ("motorcycles") -> 3 items (2500, 2500, 5000 km)
	moto, err := carRepo.Create(ctx, domaincar.Car{
		OwnerID:        "user-1",
		Name:           "Yamaha Fazer",
		Manufacturer:   "Yamaha",
		Model:          "FZ25",
		InitialMileage: 5000,
		LastMileage:    7600, // 2600 km since baseline -> oil & oil filter overdue!
		VehicleType:    "motorcycles",
		CreatedAt:      fixedNow.AddDate(0, -3, 0),
	})
	if err != nil {
		t.Fatalf("failed to create motorcycle: %v", err)
	}

	upcomingMoto, err := svc.GetUpcomingByCar(ctx, "user-1", moto.ID)
	if err != nil {
		t.Fatalf("unexpected error for motorcycle: %v", err)
	}
	if len(upcomingMoto.Items) != 3 {
		t.Fatalf("expected 3 periodic items for motorcycle, got %d", len(upcomingMoto.Items))
	}
	if upcomingMoto.Items[0].IntervalKm != 2500 {
		t.Errorf("expected 2500 km oil interval for motorcycle, got %d", upcomingMoto.Items[0].IntervalKm)
	}
	if upcomingMoto.Items[0].Status != domainmaintenance.StatusOverdue {
		t.Errorf("expected motorcycle oil status overdue, got %q", upcomingMoto.Items[0].Status)
	}
	if upcomingMoto.Items[2].IntervalKm != 5000 || upcomingMoto.Items[2].Status != domainmaintenance.StatusOK {
		t.Errorf("expected motorcycle air filter (5000 km) to be ok, got interval=%d status=%q", upcomingMoto.Items[2].IntervalKm, upcomingMoto.Items[2].Status)
	}

	// 3. Test Fleet Summary & Authorization
	summary, err := svc.GetUpcomingSummary(ctx, "user-1")
	if err != nil {
		t.Fatalf("unexpected error in summary: %v", err)
	}
	if summary.TotalVehicles != 2 {
		t.Errorf("expected 2 vehicles in summary, got %d", summary.TotalVehicles)
	}
	if summary.Status != domainmaintenance.StatusOverdue {
		t.Errorf("expected fleet status overdue, got %q", summary.Status)
	}

	// Unauthorized user cannot access sedan upcoming maintenances
	if _, err := svc.GetUpcomingByCar(ctx, "other-user", sedan.ID); !errors.Is(err, ErrMaintenanceForbidden) {
		t.Errorf("expected ErrMaintenanceForbidden for unauthorized user, got %v", err)
	}
}

func TestService_GetUpcomingByCar_CustomIntervals(t *testing.T) {
	ctx := context.Background()
	carRepo := carMemory.NewRepository()
	maintRepo := maintMemory.NewRepository()
	svc := NewService(maintRepo, carRepo)

	fixedNow := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return fixedNow }

	// Vehicle with custom intervals: Oil at 5,000 km / 6 months; Fuel filter at 15,000 km / 24 months
	car, err := carRepo.Create(ctx, domaincar.Car{
		OwnerID:        "user-1",
		Name:           "Corollinha Custom",
		Manufacturer:   "Toyota",
		Model:          "Corolla",
		InitialMileage: 50000,
		LastMileage:    54200, // 4,200 km since baseline
		VehicleType:    "cars",
		CreatedAt:      fixedNow.AddDate(0, -5, 0), // 5 months ago
		MaintenanceIntervals: map[string]domaincar.MaintenanceIntervalConfig{
			domainmaintenance.TypeEngineOil: {
				IntervalKm:     5000,
				IntervalMonths: 6,
			},
			domainmaintenance.TypeFuelFilter: {
				IntervalKm:     15000,
				IntervalMonths: 24,
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create car: %v", err)
	}

	upcoming, err := svc.GetUpcomingByCar(ctx, "user-1", car.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 1. Engine Oil: custom 5000 km (next at 55000 -> 800 km left -> Upcoming!)
	var oilItem domainmaintenance.UpcomingMaintenanceItem
	for _, it := range upcoming.Items {
		if it.Type == domainmaintenance.TypeEngineOil {
			oilItem = it
		}
	}
	if oilItem.IntervalKm != 5000 {
		t.Errorf("expected custom intervalKm 5000, got %d", oilItem.IntervalKm)
	}
	if oilItem.IntervalMonths != 6 {
		t.Errorf("expected custom intervalMonths 6, got %d", oilItem.IntervalMonths)
	}
	if oilItem.NextDueMileage != 55000 {
		t.Errorf("expected nextDueMileage 55000, got %d", oilItem.NextDueMileage)
	}
	if oilItem.RemainingKm != 800 {
		t.Errorf("expected 800 km remaining, got %d", oilItem.RemainingKm)
	}
	if oilItem.Status != domainmaintenance.StatusUpcoming {
		t.Errorf("expected status upcoming, got %q", oilItem.Status)
	}

	// 2. Oil Filter: not customized, should use default (10000 km, 12 months -> next at 60000 -> 5800 km left -> OK)
	var filterItem domainmaintenance.UpcomingMaintenanceItem
	for _, it := range upcoming.Items {
		if it.Type == domainmaintenance.TypeOilFilter {
			filterItem = it
		}
	}
	if filterItem.IntervalKm != 10000 {
		t.Errorf("expected default intervalKm 10000, got %d", filterItem.IntervalKm)
	}
	if filterItem.IntervalMonths != 12 {
		t.Errorf("expected default intervalMonths 12, got %d", filterItem.IntervalMonths)
	}
	if filterItem.Status != domainmaintenance.StatusOK {
		t.Errorf("expected default filter to be ok, got %q", filterItem.Status)
	}

	// 3. Fuel Filter: custom 15000 km, 24 months
	var fuelFilterItem domainmaintenance.UpcomingMaintenanceItem
	for _, it := range upcoming.Items {
		if it.Type == domainmaintenance.TypeFuelFilter {
			fuelFilterItem = it
		}
	}
	if fuelFilterItem.IntervalKm != 15000 {
		t.Errorf("expected custom intervalKm 15000, got %d", fuelFilterItem.IntervalKm)
	}
	if fuelFilterItem.IntervalMonths != 24 {
		t.Errorf("expected custom intervalMonths 24, got %d", fuelFilterItem.IntervalMonths)
	}
}
