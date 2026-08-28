package car

import (
	"context"
	"errors"
	"testing"

	"cvmc/internal/application/usecase/auth"
	"cvmc/internal/infrastructure/auth/bcrypt"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	carrepo "cvmc/internal/infrastructure/car/memory"
	memoryuser "cvmc/internal/infrastructure/user/memory"
)

func TestCarServiceCreateRequiresVerifiedEmail(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	authService := auth.NewService(users, hasher, tokens, nil)
	ctx := context.Background()

	// Register user (unverified by default)
	unverifiedUser, err := authService.Register(ctx, auth.RegisterInput{Name: "Unverified", Email: "unverified@example.com", Password: "secret123"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	service := NewService(carrepo.NewRepository(), users)

	// Attempt to create car with unverified email
	_, err = service.Create(ctx, unverifiedUser.User.ID, CreateInput{
		Name:            "Carro 1",
		Manufacturer:    "Fiat",
		Model:           "Uno",
		YearManufacture: 2010,
		YearModel:       2011,
		LastMileage:     10000,
	})
	if !errors.Is(err, ErrEmailNotVerified) {
		t.Fatalf("expected ErrEmailNotVerified, got: %v", err)
	}

	// Now mark user as verified
	unverifiedUser.User.EmailVerified = true
	_, err = users.Update(ctx, unverifiedUser.User)
	if err != nil {
		t.Fatalf("update user failed: %v", err)
	}

	// Should now succeed
	created, err := service.Create(ctx, unverifiedUser.User.ID, CreateInput{
		Name:            "Carro 1",
		Manufacturer:    "Fiat",
		Model:           "Uno",
		YearManufacture: 2010,
		YearModel:       2011,
		LastMileage:     10000,
	})
	if err != nil {
		t.Fatalf("expected create to succeed after verification, got: %v", err)
	}
	if created.Name != "Carro 1" {
		t.Fatalf("unexpected car name: %s", created.Name)
	}
}

func TestCarServiceVehicleLimit(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	authService := auth.NewService(users, hasher, tokens, nil)
	ctx := context.Background()

	registered, err := authService.Register(ctx, auth.RegisterInput{Name: "Collector", Email: "collector@example.com", Password: "secret123"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Mark verified, default maxVehicles = 3
	registered.User.EmailVerified = true
	registered.User.MaxVehicles = 3
	_, err = users.Update(ctx, registered.User)
	if err != nil {
		t.Fatalf("update user failed: %v", err)
	}

	service := NewService(carrepo.NewRepository(), users)

	// Create 3 cars (should all succeed)
	for i := 1; i <= 3; i++ {
		_, err := service.Create(ctx, registered.User.ID, CreateInput{
			Name:            "Carro",
			Manufacturer:    "Fiat",
			Model:           "Uno",
			YearManufacture: 2010,
			YearModel:       2011,
			LastMileage:     10000,
		})
		if err != nil {
			t.Fatalf("failed to create car %d: %v", i, err)
		}
	}

	// 4th car should fail with ErrVehicleLimitReached
	_, err = service.Create(ctx, registered.User.ID, CreateInput{
		Name:            "Carro 4",
		Manufacturer:    "Fiat",
		Model:           "Uno",
		YearManufacture: 2010,
		YearModel:       2011,
		LastMileage:     10000,
	})
	if !errors.Is(err, ErrVehicleLimitReached) {
		t.Fatalf("expected ErrVehicleLimitReached for 4th car, got: %v", err)
	}

	// Increase user limit to 5
	registered.User.MaxVehicles = 5
	_, err = users.Update(ctx, registered.User)
	if err != nil {
		t.Fatalf("update user failed: %v", err)
	}

	// Now 4th car should succeed
	_, err = service.Create(ctx, registered.User.ID, CreateInput{
		Name:            "Carro 4",
		Manufacturer:    "Fiat",
		Model:           "Uno",
		YearManufacture: 2010,
		YearModel:       2011,
		LastMileage:     10000,
	})
	if err != nil {
		t.Fatalf("expected create 4th car to succeed after limit increase, got: %v", err)
	}
}

func TestCarServiceShareAndList(t *testing.T) {
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
	_, _ = users.Update(ctx, owner.User)

	viewer, err := authService.Register(ctx, auth.RegisterInput{Name: "Viewer", Email: "viewer@example.com", Password: "secret123"})
	if err != nil {
		t.Fatalf("viewer register failed: %v", err)
	}

	service := NewService(carrepo.NewRepository(), users)
	created, err := service.Create(ctx, owner.User.ID, CreateInput{Name: "Carro 1", Manufacturer: "Fiat", Model: "Uno", YearManufacture: 2010, YearModel: 2011, LastMileage: 10000})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if _, err := service.Share(ctx, owner.User.ID, created.ID, viewer.User.Email); err != nil {
		t.Fatalf("share failed: %v", err)
	}

	items, err := service.List(ctx, viewer.User.ID)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 shared car, got %d", len(items))
	}
}

func TestCarServiceUpdate(t *testing.T) {
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
	_, _ = users.Update(ctx, owner.User)

	otherUser, err := authService.Register(ctx, auth.RegisterInput{Name: "Other", Email: "other@example.com", Password: "secret123"})
	if err != nil {
		t.Fatalf("other register failed: %v", err)
	}

	service := NewService(carrepo.NewRepository(), users)
	created, err := service.Create(ctx, owner.User.ID, CreateInput{
		Name:            "Meu Carro",
		Manufacturer:    "Honda",
		Model:           "Civic",
		YearManufacture: 2020,
		YearModel:       2021,
		LastMileage:     30000,
		ImageUrl:        "https://example.com/car.jpg",
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if created.ImageUrl != "https://example.com/car.jpg" {
		t.Fatalf("expected ImageUrl to be set, got %s", created.ImageUrl)
	}

	// Successful update
	updated, err := service.Update(ctx, owner.User.ID, created.ID, UpdateInput{
		Name:            "Civic Novo",
		LastMileage:     35000,
		YearManufacture: 2020,
		YearModel:       2021,
		ImageUrl:        "https://example.com/car-updated.jpg",
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Name != "Civic Novo" || updated.LastMileage != 35000 || updated.Manufacturer != "Honda" || updated.ImageUrl != "https://example.com/car-updated.jpg" {
		t.Fatalf("unexpected updated car values: %+v", updated)
	}

	// Forbidden update by non-owner
	_, err = service.Update(ctx, otherUser.User.ID, created.ID, UpdateInput{Name: "Hacked"})
	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got: %v", err)
	}

	// Invalid mileage
	_, err = service.Update(ctx, owner.User.ID, created.ID, UpdateInput{LastMileage: -10})
	if err != ErrInvalidPayload {
		t.Fatalf("expected ErrInvalidPayload for negative mileage, got: %v", err)
	}
}
