package fipe

import (
	"context"
	"errors"
	"testing"
	"time"

	domainfipe "cvmc/internal/domain/fipe"
	memoryfipe "cvmc/internal/infrastructure/fipe/memory"
)

type mockExternalClient struct {
	brandsFunc func(ctx context.Context, vt domainfipe.VehicleType) ([]domainfipe.Brand, error)
	modelsFunc func(ctx context.Context, vt domainfipe.VehicleType, brandCode string) ([]domainfipe.Model, error)
	yearsFunc  func(ctx context.Context, vt domainfipe.VehicleType, brandCode, modelCode string) ([]domainfipe.Year, error)
	detailFunc func(ctx context.Context, vt domainfipe.VehicleType, brandCode, modelCode, yearCode string) (domainfipe.VehicleDetail, error)
}

func (m *mockExternalClient) FetchBrands(ctx context.Context, vt domainfipe.VehicleType) ([]domainfipe.Brand, error) {
	if m.brandsFunc != nil {
		return m.brandsFunc(ctx, vt)
	}
	return nil, errors.New("not implemented")
}

func (m *mockExternalClient) FetchModels(ctx context.Context, vt domainfipe.VehicleType, brandCode string) ([]domainfipe.Model, error) {
	if m.modelsFunc != nil {
		return m.modelsFunc(ctx, vt, brandCode)
	}
	return nil, errors.New("not implemented")
}

func (m *mockExternalClient) FetchYears(ctx context.Context, vt domainfipe.VehicleType, brandCode, modelCode string) ([]domainfipe.Year, error) {
	if m.yearsFunc != nil {
		return m.yearsFunc(ctx, vt, brandCode, modelCode)
	}
	return nil, errors.New("not implemented")
}

func (m *mockExternalClient) FetchVehicleDetail(ctx context.Context, vt domainfipe.VehicleType, brandCode, modelCode, yearCode string) (domainfipe.VehicleDetail, error) {
	if m.detailFunc != nil {
		return m.detailFunc(ctx, vt, brandCode, modelCode, yearCode)
	}
	return domainfipe.VehicleDetail{}, errors.New("not implemented")
}

func TestService_GetBrands_CacheAndFallback(t *testing.T) {
	ctx := context.Background()
	repo := memoryfipe.NewRepository()
	fetchCount := 0

	client := &mockExternalClient{
		brandsFunc: func(ctx context.Context, vt domainfipe.VehicleType) ([]domainfipe.Brand, error) {
			fetchCount++
			if fetchCount == 2 {
				return nil, errors.New("fipe api 500 error")
			}
			return []domainfipe.Brand{
				{Code: "59", Name: "VW - VolksWagen"},
				{Code: "21", Name: "Fiat"},
			}, nil
		},
	}

	currentTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	svc := NewService(repo, client)
	svc.now = func() time.Time { return currentTime }

	// 1. Initial call: cache miss -> fetch from API
	brands, err := svc.GetBrands(ctx, domainfipe.VehicleTypeCars)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(brands) != 2 {
		t.Fatalf("expected 2 brands, got %d", len(brands))
	}
	if fetchCount != 1 {
		t.Fatalf("expected fetchCount 1, got %d", fetchCount)
	}

	// 2. Immediate second call: cache hit -> no fetch
	brands, err = svc.GetBrands(ctx, domainfipe.VehicleTypeCars)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(brands) != 2 {
		t.Fatalf("expected 2 brands, got %d", len(brands))
	}
	if fetchCount != 1 {
		t.Fatalf("expected fetchCount still 1, got %d", fetchCount)
	}

	// 3. Fast forward past TTL (91 days) + API fails -> should return stale data
	currentTime = currentTime.Add(91 * 24 * time.Hour)
	brands, err = svc.GetBrands(ctx, domainfipe.VehicleTypeCars)
	if err != nil {
		t.Fatalf("expected no error from stale fallback, got %v", err)
	}
	if len(brands) != 2 {
		t.Fatalf("expected 2 brands from fallback, got %d", len(brands))
	}
	if fetchCount != 2 {
		t.Fatalf("expected fetchCount 2, got %d", fetchCount)
	}
}

func TestService_GetModelsAndYearsAndDetails(t *testing.T) {
	ctx := context.Background()
	repo := memoryfipe.NewRepository()

	modelsFetchCount := 0
	yearsFetchCount := 0
	detailFetchCount := 0

	client := &mockExternalClient{
		brandsFunc: func(ctx context.Context, vt domainfipe.VehicleType) ([]domainfipe.Brand, error) {
			return []domainfipe.Brand{{Code: "59", Name: "VW - VolksWagen"}}, nil
		},
		modelsFunc: func(ctx context.Context, vt domainfipe.VehicleType, brandCode string) ([]domainfipe.Model, error) {
			modelsFetchCount++
			return []domainfipe.Model{{Code: "5940", Name: "Gol 1.0"}}, nil
		},
		yearsFunc: func(ctx context.Context, vt domainfipe.VehicleType, brandCode, modelCode string) ([]domainfipe.Year, error) {
			yearsFetchCount++
			return []domainfipe.Year{{Code: "2023-1", Name: "2023 Gasolina"}}, nil
		},
		detailFunc: func(ctx context.Context, vt domainfipe.VehicleType, brandCode, modelCode, yearCode string) (domainfipe.VehicleDetail, error) {
			detailFetchCount++
			return domainfipe.VehicleDetail{
				Brand:          "VW - VolksWagen",
				Model:          "Gol 1.0",
				ModelYear:      2023,
				Price:          "R$ 55.400,00",
				PriceValue:     55400,
				CodeFipe:       "005487-9",
				Fuel:           "Gasolina",
				ReferenceMonth: "agosto de 2026",
				VehicleType:    1,
			}, nil
		},
	}

	svc := NewService(repo, client)

	// Brands
	brands, err := svc.GetBrands(ctx, domainfipe.VehicleTypeCars)
	if err != nil || len(brands) == 0 {
		t.Fatalf("failed to get brands: %v", err)
	}

	// Models - Call 1 (miss)
	models, err := svc.GetModels(ctx, domainfipe.VehicleTypeCars, "59")
	if err != nil || len(models) == 0 {
		t.Fatalf("failed to get models: %v", err)
	}
	if models[0].Name != "Gol 1.0" {
		t.Fatalf("expected Gol 1.0, got %s", models[0].Name)
	}
	if modelsFetchCount != 1 {
		t.Fatalf("expected modelsFetchCount 1, got %d", modelsFetchCount)
	}

	// Models - Call 2 (cache hit, no external call)
	models, err = svc.GetModels(ctx, domainfipe.VehicleTypeCars, "59")
	if err != nil || len(models) == 0 {
		t.Fatalf("failed to get models from cache: %v", err)
	}
	if modelsFetchCount != 1 {
		t.Fatalf("expected modelsFetchCount still 1 on cache hit, got %d", modelsFetchCount)
	}

	// Years - Call 1 (miss)
	years, err := svc.GetYears(ctx, domainfipe.VehicleTypeCars, "59", "5940")
	if err != nil || len(years) == 0 {
		t.Fatalf("failed to get years: %v", err)
	}
	if years[0].Code != "2023-1" {
		t.Fatalf("expected 2023-1, got %s", years[0].Code)
	}
	if yearsFetchCount != 1 {
		t.Fatalf("expected yearsFetchCount 1, got %d", yearsFetchCount)
	}

	// Years - Call 2 (cache hit)
	years, err = svc.GetYears(ctx, domainfipe.VehicleTypeCars, "59", "5940")
	if err != nil || len(years) == 0 {
		t.Fatalf("failed to get years from cache: %v", err)
	}
	if yearsFetchCount != 1 {
		t.Fatalf("expected yearsFetchCount still 1 on cache hit, got %d", yearsFetchCount)
	}

	// Detail - Call 1 (miss)
	detail, err := svc.GetVehicleDetail(ctx, domainfipe.VehicleTypeCars, "59", "5940", "2023-1")
	if err != nil {
		t.Fatalf("failed to get vehicle detail: %v", err)
	}
	if detail.CodeFipe != "005487-9" {
		t.Fatalf("expected 005487-9, got %s", detail.CodeFipe)
	}
	if detailFetchCount != 1 {
		t.Fatalf("expected detailFetchCount 1, got %d", detailFetchCount)
	}

	// Detail - Call 2 (cache hit)
	detail, err = svc.GetVehicleDetail(ctx, domainfipe.VehicleTypeCars, "59", "5940", "2023-1")
	if err != nil {
		t.Fatalf("failed to get vehicle detail from cache: %v", err)
	}
	if detail.CodeFipe != "005487-9" {
		t.Fatalf("expected 005487-9, got %s", detail.CodeFipe)
	}
	if detailFetchCount != 1 {
		t.Fatalf("expected detailFetchCount still 1 on cache hit, got %d", detailFetchCount)
	}
}
