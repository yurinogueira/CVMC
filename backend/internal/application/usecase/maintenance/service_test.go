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

func TestMaintenanceServiceUpdateAndForbidden(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	authService := auth.NewService(users, hasher, tokens, nil)
	ctx := context.Background()

	owner, err := authService.Register(ctx, auth.RegisterInput{Name: "Owner", Email: "owner2@example.com", Password: "secret123"})
	if err != nil {
		t.Fatalf("owner register failed: %v", err)
	}
	other, err := authService.Register(ctx, auth.RegisterInput{Name: "Other", Email: "other2@example.com", Password: "secret123"})
	if err != nil {
		t.Fatalf("other register failed: %v", err)
	}

	owner.User.EmailVerified = true
	other.User.EmailVerified = true
	_, _ = users.Update(ctx, owner.User)
	_, _ = users.Update(ctx, other.User)

	cars := carrepo.NewRepository()
	carService := car.NewService(cars, users)
	createdCar, err := carService.Create(ctx, owner.User.ID, car.CreateInput{Name: "Carro 2", Manufacturer: "VW", Model: "Gol", YearManufacture: 2015, YearModel: 2016, LastMileage: 50000})
	if err != nil {
		t.Fatalf("create car failed: %v", err)
	}

	service := NewService(maintrepo.NewRepository(), cars)
	created, err := service.Create(ctx, owner.User.ID, createdCar.ID, CreateInput{
		Title:   "Troca de óleo inicial",
		Date:    time.Now().UTC(),
		Mileage: 51000,
	})
	if err != nil {
		t.Fatalf("create maintenance failed: %v", err)
	}

	// 1. Successful update of fields and mileage > car.LastMileage
	costVal := 450.0
	updated, err := service.Update(ctx, owner.User.ID, created.ID, UpdateInput{
		Title:       "Revisão dos 55.000 km",
		Description: "Troca de óleo e filtro de combustível",
		Date:        time.Now().UTC(),
		Mileage:     55000,
		Types:       []string{"Óleo de Motor", "Filtro de Combustível"},
		Cost:        &costVal,
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Title != "Revisão dos 55.000 km" || updated.Mileage != 55000 || updated.Cost == nil || *updated.Cost != 450.0 {
		t.Fatalf("unexpected updated fields: %+v", updated)
	}

	// Verify car mileage was bumped to 55000
	c, err := cars.GetByID(ctx, createdCar.ID)
	if err != nil || c.LastMileage != 55000 {
		t.Fatalf("expected car last mileage 55000, got %d", c.LastMileage)
	}

	// 2. Forbidden when other user tries to update
	_, err = service.Update(ctx, other.User.ID, created.ID, UpdateInput{Title: "Hacked"})
	if err != ErrMaintenanceForbidden {
		t.Fatalf("expected ErrMaintenanceForbidden, got %v", err)
	}

	// 3. Invalid when title is empty
	_, err = service.Update(ctx, owner.User.ID, created.ID, UpdateInput{Title: ""})
	if err != ErrMaintenanceInvalid {
		t.Fatalf("expected ErrMaintenanceInvalid, got %v", err)
	}

	// 4. Not found when maintenance does not exist
	_, err = service.Update(ctx, owner.User.ID, "non-existent-id", UpdateInput{Title: "Alguma coisa"})
	if err != ErrMaintenanceNotFound {
		t.Fatalf("expected ErrMaintenanceNotFound, got %v", err)
	}
}

func TestMaintenanceServiceDeleteAndForbidden(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	authService := auth.NewService(users, hasher, tokens, nil)
	ctx := context.Background()

	owner, _ := authService.Register(ctx, auth.RegisterInput{Name: "Owner", Email: "owner3@example.com", Password: "secret123"})
	other, _ := authService.Register(ctx, auth.RegisterInput{Name: "Other", Email: "other3@example.com", Password: "secret123"})

	owner.User.EmailVerified = true
	other.User.EmailVerified = true
	_, _ = users.Update(ctx, owner.User)
	_, _ = users.Update(ctx, other.User)

	cars := carrepo.NewRepository()
	carService := car.NewService(cars, users)
	createdCar, _ := carService.Create(ctx, owner.User.ID, car.CreateInput{Name: "Carro 3", Manufacturer: "Honda", Model: "Fit", YearManufacture: 2018, YearModel: 2019, LastMileage: 60000})

	service := NewService(maintrepo.NewRepository(), cars)
	created, err := service.Create(ctx, owner.User.ID, createdCar.ID, CreateInput{
		Title:   "Troca de pastilhas",
		Date:    time.Now().UTC(),
		Mileage: 60000,
	})
	if err != nil {
		t.Fatalf("create maintenance failed: %v", err)
	}

	// 1. Forbidden when other user tries to delete
	err = service.Delete(ctx, other.User.ID, created.ID)
	if err != ErrMaintenanceForbidden {
		t.Fatalf("expected ErrMaintenanceForbidden, got %v", err)
	}

	// 2. Successful delete by owner
	err = service.Delete(ctx, owner.User.ID, created.ID)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// 3. Not found when trying to delete again
	err = service.Delete(ctx, owner.User.ID, created.ID)
	if err != ErrMaintenanceNotFound {
		t.Fatalf("expected ErrMaintenanceNotFound, got %v", err)
	}
}

func TestMaintenanceServiceGet(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	authService := auth.NewService(users, hasher, tokens, nil)
	ctx := context.Background()

	owner, _ := authService.Register(ctx, auth.RegisterInput{Name: "Owner", Email: "owner4@example.com", Password: "secret123"})
	other, _ := authService.Register(ctx, auth.RegisterInput{Name: "Other", Email: "other4@example.com", Password: "secret123"})

	owner.User.EmailVerified = true
	other.User.EmailVerified = true
	_, _ = users.Update(ctx, owner.User)
	_, _ = users.Update(ctx, other.User)

	cars := carrepo.NewRepository()
	carService := car.NewService(cars, users)
	createdCar, _ := carService.Create(ctx, owner.User.ID, car.CreateInput{Name: "Carro 4", Manufacturer: "Toyota", Model: "Etios", YearManufacture: 2017, YearModel: 2018, LastMileage: 80000})

	service := NewService(maintrepo.NewRepository(), cars)
	created, err := service.Create(ctx, owner.User.ID, createdCar.ID, CreateInput{
		Title:   "Alinhamento e Balanceamento",
		Date:    time.Now().UTC(),
		Mileage: 80000,
	})
	if err != nil {
		t.Fatalf("create maintenance failed: %v", err)
	}

	// 1. Get success by owner
	item, err := service.Get(ctx, owner.User.ID, created.ID)
	if err != nil || item.ID != created.ID {
		t.Fatalf("expected item, got err: %v", err)
	}

	// 2. Forbidden when other user tries to get
	_, err = service.Get(ctx, other.User.ID, created.ID)
	if err != ErrMaintenanceForbidden {
		t.Fatalf("expected ErrMaintenanceForbidden, got %v", err)
	}

	// 3. Not found when ID does not exist
	_, err = service.Get(ctx, owner.User.ID, "non-existent")
	if err != ErrMaintenanceNotFound {
		t.Fatalf("expected ErrMaintenanceNotFound, got %v", err)
	}
}
