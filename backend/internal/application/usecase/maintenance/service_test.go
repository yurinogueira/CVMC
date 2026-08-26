package maintenance

import (
	"context"
	"testing"
	"time"

	"cvmc/internal/application/usecase/auth"
	"cvmc/internal/application/usecase/car"
	"cvmc/internal/infrastructure/auth/bcrypt"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	carrepo "cvmc/internal/infrastructure/car/memory"
	maintrepo "cvmc/internal/infrastructure/maintenance/memory"
	memoryuser "cvmc/internal/infrastructure/user/memory"
)

func TestMaintenanceServiceUpdatesMileage(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
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
	createdCar, err := carService.Create(ctx, owner.User.ID, car.CreateInput{Name: "Carro 1", Manufacturer: "Fiat", Model: "Uno", YearManufacture: 2010, YearModel: 2011, LastMileage: 10000})
	if err != nil {
		t.Fatalf("create car failed: %v", err)
	}

	service := NewService(maintrepo.NewRepository(), cars)
	createdMaintenance, err := service.Create(ctx, owner.User.ID, createdCar.ID, CreateInput{Title: "Troca de óleo", Description: "Filtro e óleo", Date: time.Now().UTC(), Mileage: 12000})
	if err != nil {
		t.Fatalf("create maintenance failed: %v", err)
	}
	if createdMaintenance.Mileage != 12000 {
		t.Fatalf("unexpected mileage: %d", createdMaintenance.Mileage)
	}
	updatedCar, err := cars.GetByID(ctx, createdCar.ID)
	if err != nil {
		t.Fatalf("get car failed: %v", err)
	}
	if updatedCar.LastMileage != 12000 {
		t.Fatalf("expected last mileage updated, got %d", updatedCar.LastMileage)
	}
}
