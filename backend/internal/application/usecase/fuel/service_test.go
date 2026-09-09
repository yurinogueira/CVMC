package fuel

import (
	"context"
	"testing"
	"time"

	"cvmc/internal/application/usecase/auth"
	"cvmc/internal/application/usecase/car"
	"cvmc/internal/infrastructure/auth/bcrypt"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	carrepo "cvmc/internal/infrastructure/car/memory"
	fuelrepo "cvmc/internal/infrastructure/fuel/memory"
	memoryuser "cvmc/internal/infrastructure/user/memory"
)

func setupTest(t *testing.T) (*Service, string, string) {
	t.Helper()
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret-32-chars-minimum-ok", "refresh-secret-32-chars-minimum-ok")
	authService := auth.NewService(users, hasher, tokens, nil)
	ctx := context.Background()

	owner, err := authService.Register(ctx, auth.RegisterInput{Name: "Owner", Email: "owner@example.com", Password: "secret123"})
	if err != nil {
		t.Fatalf("owner register failed: %v", err)
	}
	owner.User.EmailVerified = true
	if _, err := users.Update(ctx, owner.User); err != nil {
		t.Fatalf("update user failed: %v", err)
	}

	cars := carrepo.NewRepository()
	carService := car.NewService(cars, users)
	createdCar, err := carService.Create(ctx, owner.User.ID, car.CreateInput{
		Name:            "Carro 1",
		Manufacturer:    "Fiat",
		Model:           "Uno",
		YearManufacture: 2015,
		YearModel:       2016,
		LastMileage:     10000,
	})
	if err != nil {
		t.Fatalf("create car failed: %v", err)
	}

	service := NewService(fuelrepo.NewRepository(), cars)
	return service, owner.User.ID, createdCar.ID
}

func TestFuelServiceCalculatesTotalAndPricePerLiter(t *testing.T) {
	service, ownerID, carID := setupTest(t)
	ctx := context.Background()

	// 1. Given Liters and PricePerLiter -> calculates TotalCost
	created1, err := service.Create(ctx, ownerID, carID, CreateInput{
		Date:          time.Now().UTC(),
		FuelType:      "Gasolina Comum",
		Liters:        40.0,
		PricePerLiter: 5.50,
		TotalCost:     0,
		IsFullTank:    true,
		GasStation:    "Posto Shell",
	})
	if err != nil {
		t.Fatalf("create fueling 1 failed: %v", err)
	}
	if created1.TotalCost != 220.0 {
		t.Errorf("expected totalCost 220.0, got %f", created1.TotalCost)
	}

	// 2. Given Liters and TotalCost -> calculates PricePerLiter automatically!
	created2, err := service.Create(ctx, ownerID, carID, CreateInput{
		Date:          time.Now().UTC(),
		FuelType:      "Etanol",
		Liters:        50.0,
		PricePerLiter: 0,
		TotalCost:     200.0,
		IsFullTank:    true,
	})
	if err != nil {
		t.Fatalf("create fueling 2 failed: %v", err)
	}
	if created2.PricePerLiter != 4.0 {
		t.Errorf("expected pricePerLiter 4.0, got %f", created2.PricePerLiter)
	}
}

func TestFuelServiceValidation(t *testing.T) {
	service, ownerID, carID := setupTest(t)
	ctx := context.Background()

	// Empty fuel type
	_, err := service.Create(ctx, ownerID, carID, CreateInput{
		FuelType:      "",
		Liters:        10,
		PricePerLiter: 5,
	})
	if err != ErrFuelingInvalid {
		t.Errorf("expected ErrFuelingInvalid for empty fuel type, got %v", err)
	}

	// 0 Liters
	_, err = service.Create(ctx, ownerID, carID, CreateInput{
		FuelType:      "Gasolina",
		Liters:        0,
		PricePerLiter: 5,
	})
	if err != ErrFuelingInvalid {
		t.Errorf("expected ErrFuelingInvalid for 0 liters, got %v", err)
	}

	// No price or total
	_, err = service.Create(ctx, ownerID, carID, CreateInput{
		FuelType: "Gasolina",
		Liters:   10,
	})
	if err != ErrFuelingInvalid {
		t.Errorf("expected ErrFuelingInvalid for missing price/total, got %v", err)
	}
}

func TestFuelServiceForbiddenAndNotFound(t *testing.T) {
	service, ownerID, carID := setupTest(t)
	ctx := context.Background()

	// Non-owner create
	_, err := service.Create(ctx, "intruder-id", carID, CreateInput{
		FuelType:      "Gasolina",
		Liters:        20,
		PricePerLiter: 5,
	})
	if err != ErrFuelingForbidden {
		t.Errorf("expected ErrFuelingForbidden, got %v", err)
	}

	// Car not found
	_, err = service.Create(ctx, ownerID, "invalid-car-id", CreateInput{
		FuelType:      "Gasolina",
		Liters:        20,
		PricePerLiter: 5,
	})
	if err != ErrFuelingNotFound {
		t.Errorf("expected ErrFuelingNotFound, got %v", err)
	}
}

func TestFuelServiceUpdateAndDelete(t *testing.T) {
	service, ownerID, carID := setupTest(t)
	ctx := context.Background()

	created, err := service.Create(ctx, ownerID, carID, CreateInput{
		FuelType:      "Gasolina",
		Liters:        30,
		PricePerLiter: 5.0,
		TotalCost:     150.0,
	})
	if err != nil {
		t.Fatalf("create fueling failed: %v", err)
	}

	// Update by non-owner
	_, err = service.Update(ctx, "stranger", created.ID, UpdateInput{
		FuelType:      "Etanol",
		Liters:        35,
		PricePerLiter: 4.0,
	})
	if err != ErrFuelingForbidden {
		t.Errorf("expected ErrFuelingForbidden, got %v", err)
	}

	// Update by owner: Liters=40, TotalCost=160 -> pricePerLiter calculated to 4.0
	updated, err := service.Update(ctx, ownerID, created.ID, UpdateInput{
		FuelType:   "Etanol",
		Liters:     40,
		TotalCost:  160.0,
		GasStation: "Posto Ipiranga",
	})
	if err != nil {
		t.Fatalf("update fueling failed: %v", err)
	}
	if updated.FuelType != "Etanol" || updated.PricePerLiter != 4.0 || updated.GasStation != "Posto Ipiranga" {
		t.Errorf("unexpected updated data: %+v", updated)
	}

	// Delete by stranger
	if err := service.Delete(ctx, "stranger", created.ID); err != ErrFuelingForbidden {
		t.Errorf("expected ErrFuelingForbidden, got %v", err)
	}

	// Delete by owner
	if err := service.Delete(ctx, ownerID, created.ID); err != nil {
		t.Fatalf("delete fueling failed: %v", err)
	}

	// Get after delete
	_, err = service.Get(ctx, ownerID, created.ID)
	if err != ErrFuelingNotFound {
		t.Errorf("expected ErrFuelingNotFound after delete, got %v", err)
	}
}

func TestFuelServiceListSortedByDate(t *testing.T) {
	service, ownerID, carID := setupTest(t)
	ctx := context.Background()

	now := time.Now().UTC()
	_, _ = service.Create(ctx, ownerID, carID, CreateInput{
		Date:          now.Add(-48 * time.Hour),
		FuelType:      "Gasolina",
		Liters:        20,
		PricePerLiter: 5.0,
	})
	_, _ = service.Create(ctx, ownerID, carID, CreateInput{
		Date:          now,
		FuelType:      "Gasolina",
		Liters:        25,
		PricePerLiter: 5.2,
	})
	_, _ = service.Create(ctx, ownerID, carID, CreateInput{
		Date:          now.Add(-24 * time.Hour),
		FuelType:      "Gasolina",
		Liters:        22,
		PricePerLiter: 5.1,
	})

	list, err := service.List(ctx, ownerID, carID)
	if err != nil {
		t.Fatalf("list fuelings failed: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 items, got %d", len(list))
	}
	if !list[0].Date.After(list[1].Date) || !list[1].Date.After(list[2].Date) {
		t.Errorf("expected list to be sorted descending by date")
	}
}
