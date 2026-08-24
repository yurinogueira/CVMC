package fipe

import (
	"context"
	"time"

	domainfipe "cvmc/internal/domain/fipe"
)

type Repository interface {
	GetBrands(ctx context.Context, vehicleType domainfipe.VehicleType) ([]domainfipe.Brand, time.Time, error)
	UpsertBrands(ctx context.Context, vehicleType domainfipe.VehicleType, brands []domainfipe.Brand, syncTime time.Time) error
	GetBrandWithModels(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string) (*domainfipe.BrandDocument, error)
	UpdateModels(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string, models []domainfipe.Model, syncTime time.Time) error
	UpdateModelYears(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string, modelCode string, years []domainfipe.Year, syncTime time.Time) error
	UpdateYearDetail(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string, modelCode string, yearCode string, detail domainfipe.VehicleDetail, syncTime time.Time) error
}

type ExternalClient interface {
	FetchBrands(ctx context.Context, vehicleType domainfipe.VehicleType) ([]domainfipe.Brand, error)
	FetchModels(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string) ([]domainfipe.Model, error)
	FetchYears(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode, modelCode string) ([]domainfipe.Year, error)
	FetchVehicleDetail(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode, modelCode, yearCode string) (domainfipe.VehicleDetail, error)
}
