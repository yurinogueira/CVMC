package fipe

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	fipeport "cvmc/internal/application/ports/fipe"
	domainfipe "cvmc/internal/domain/fipe"
)

const (
	BrandsCacheTTL  = 90 * 24 * time.Hour // 3 months
	ModelsCacheTTL  = 30 * 24 * time.Hour // 1 month
	YearsCacheTTL   = 30 * 24 * time.Hour // 1 month
	DetailsCacheTTL = 7 * 24 * time.Hour  // 1 week
)

var (
	ErrFipeUnavailable = errors.New("fipe service unavailable and no cached data exists")
)

type Service struct {
	repo   fipeport.Repository
	client fipeport.ExternalClient
	now    func() time.Time
}

func NewService(repo fipeport.Repository, client fipeport.ExternalClient) *Service {
	return &Service{
		repo:   repo,
		client: client,
		now:    time.Now,
	}
}

func (s *Service) GetBrands(ctx context.Context, vehicleType domainfipe.VehicleType) ([]domainfipe.Brand, error) {
	cachedBrands, latestSync, err := s.repo.GetBrands(ctx, vehicleType)
	if err == nil && len(cachedBrands) > 0 && s.now().Sub(latestSync) < BrandsCacheTTL {
		return cachedBrands, nil
	}

	// Fetch from external API
	if s.client != nil {
		externalBrands, fetchErr := s.client.FetchBrands(ctx, vehicleType)
		if fetchErr == nil && len(externalBrands) > 0 {
			if saveErr := s.repo.UpsertBrands(ctx, vehicleType, externalBrands, s.now().UTC()); saveErr != nil {
				log.Printf("[ERROR] Failed to save brands in database: %v", saveErr)
			}
			return externalBrands, nil
		}
		log.Printf("[WARN] Failed to fetch brands from external FIPE API: %v", fetchErr)
	}

	// Stale fallback
	if len(cachedBrands) > 0 {
		log.Printf("[INFO] Returning %d stale cached brands for vehicleType=%s", len(cachedBrands), vehicleType)
		return cachedBrands, nil
	}

	return nil, ErrFipeUnavailable
}

func (s *Service) GetModels(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string) ([]domainfipe.Model, error) {
	doc, err := s.repo.GetBrandWithModels(ctx, vehicleType, brandCode)
	if err == nil && doc != nil && len(doc.Models) > 0 && s.now().Sub(doc.ModelsLastSyncAt) < ModelsCacheTTL {
		models := make([]domainfipe.Model, len(doc.Models))
		for i, m := range doc.Models {
			models[i] = domainfipe.Model{Code: m.Code, Name: m.Name}
		}
		return models, nil
	}

	// Fetch from external API
	if s.client != nil {
		externalModels, fetchErr := s.client.FetchModels(ctx, vehicleType, brandCode)
		if fetchErr == nil && len(externalModels) > 0 {
			if saveErr := s.repo.UpdateModels(ctx, vehicleType, brandCode, externalModels, s.now().UTC()); saveErr != nil {
				log.Printf("[ERROR] Failed to save models in database: %v", saveErr)
			}
			return externalModels, nil
		}
		log.Printf("[WARN] Failed to fetch models from external FIPE API: %v", fetchErr)
	}

	// Stale fallback
	if doc != nil && len(doc.Models) > 0 {
		log.Printf("[INFO] Returning %d stale cached models for brandCode=%s", len(doc.Models), brandCode)
		models := make([]domainfipe.Model, len(doc.Models))
		for i, m := range doc.Models {
			models[i] = domainfipe.Model{Code: m.Code, Name: m.Name}
		}
		return models, nil
	}

	return nil, ErrFipeUnavailable
}

func (s *Service) GetYears(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode, modelCode string) ([]domainfipe.Year, error) {
	doc, err := s.repo.GetBrandWithModels(ctx, vehicleType, brandCode)
	var existingModel *domainfipe.ModelDocument
	if err == nil && doc != nil {
		for i := range doc.Models {
			if doc.Models[i].Code == modelCode {
				existingModel = &doc.Models[i]
				break
			}
		}
	}

	if existingModel != nil && len(existingModel.Years) > 0 && s.now().Sub(existingModel.YearsLastSyncAt) < YearsCacheTTL {
		years := make([]domainfipe.Year, len(existingModel.Years))
		for i, y := range existingModel.Years {
			years[i] = domainfipe.Year{Code: y.Code, Name: y.Name}
		}
		return years, nil
	}

	// Fetch from external API
	if s.client != nil {
		externalYears, fetchErr := s.client.FetchYears(ctx, vehicleType, brandCode, modelCode)
		if fetchErr == nil && len(externalYears) > 0 {
			if saveErr := s.repo.UpdateModelYears(ctx, vehicleType, brandCode, modelCode, externalYears, s.now().UTC()); saveErr != nil {
				log.Printf("[ERROR] Failed to save years in database: %v", saveErr)
			}
			return externalYears, nil
		}
		log.Printf("[WARN] Failed to fetch years from external FIPE API: %v", fetchErr)
	}

	// Stale fallback
	if existingModel != nil && len(existingModel.Years) > 0 {
		log.Printf("[INFO] Returning %d stale cached years for modelCode=%s", len(existingModel.Years), modelCode)
		years := make([]domainfipe.Year, len(existingModel.Years))
		for i, y := range existingModel.Years {
			years[i] = domainfipe.Year{Code: y.Code, Name: y.Name}
		}
		return years, nil
	}

	return nil, ErrFipeUnavailable
}

func (s *Service) GetVehicleDetail(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode, modelCode, yearCode string) (domainfipe.VehicleDetail, error) {
	doc, err := s.repo.GetBrandWithModels(ctx, vehicleType, brandCode)
	var existingYear *domainfipe.YearDocument
	var modelName string
	if err == nil && doc != nil {
		for i := range doc.Models {
			if doc.Models[i].Code == modelCode {
				modelName = doc.Models[i].Name
				for j := range doc.Models[i].Years {
					if doc.Models[i].Years[j].Code == yearCode {
						existingYear = &doc.Models[i].Years[j]
						break
					}
				}
				break
			}
		}
	}

	if existingYear != nil && existingYear.Price != "" && s.now().Sub(existingYear.DetailsLastSyncAt) < DetailsCacheTTL {
		return toVehicleDetail(doc.Name, modelName, existingYear, vehicleType), nil
	}

	// Fetch from external API
	if s.client != nil {
		detail, fetchErr := s.client.FetchVehicleDetail(ctx, vehicleType, brandCode, modelCode, yearCode)
		if fetchErr == nil && detail.Price != "" {
			if saveErr := s.repo.UpdateYearDetail(ctx, vehicleType, brandCode, modelCode, yearCode, detail, s.now().UTC()); saveErr != nil {
				log.Printf("[ERROR] Failed to save vehicle detail in database: %v", saveErr)
			}
			return detail, nil
		}
		log.Printf("[WARN] Failed to fetch vehicle detail from external FIPE API: %v", fetchErr)
	}

	// Stale fallback
	if existingYear != nil && existingYear.Price != "" {
		log.Printf("[INFO] Returning stale cached vehicle detail for yearCode=%s", yearCode)
		return toVehicleDetail(doc.Name, modelName, existingYear, vehicleType), nil
	}

	return domainfipe.VehicleDetail{}, ErrFipeUnavailable
}

func toVehicleDetail(brandName, modelName string, y *domainfipe.YearDocument, vt domainfipe.VehicleType) domainfipe.VehicleDetail {
	yearInt := 0
	parts := strings.Split(y.Code, "-")
	if len(parts) > 0 {
		yearInt, _ = strconv.Atoi(parts[0])
	}
	vtInt := 1
	if vt == domainfipe.VehicleTypeMotorcycles {
		vtInt = 2
	} else if vt == domainfipe.VehicleTypeTrucks {
		vtInt = 3
	}

	return domainfipe.VehicleDetail{
		Brand:          brandName,
		Model:          modelName,
		ModelYear:      yearInt,
		Price:          y.Price,
		PriceValue:     y.PriceValue,
		CodeFipe:       y.FIPECode,
		Fuel:           y.Fuel,
		ReferenceMonth: y.ReferenceMonth,
		VehicleType:    vtInt,
	}
}
