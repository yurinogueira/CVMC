package car

import (
	"context"
	"testing"

	"cvmc/internal/application/usecase/auth"
	"cvmc/internal/infrastructure/auth/bcrypt"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	carrepo "cvmc/internal/infrastructure/car/memory"
	memoryuser "cvmc/internal/infrastructure/user/memory"
)

func TestCarServiceShareAndList(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	authService := auth.NewService(users, hasher, tokens)
	ctx := context.Background()

	owner, err := authService.Register(ctx, auth.RegisterInput{Name: "Owner", Email: "owner@example.com", Password: "secret123"})
	if err != nil {
		t.Fatalf("owner register failed: %v", err)
	}
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
