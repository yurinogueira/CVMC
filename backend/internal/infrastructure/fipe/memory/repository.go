package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	domainfipe "cvmc/internal/domain/fipe"
)

type Repository struct {
	mu     sync.RWMutex
	brands map[string]*domainfipe.BrandDocument // key: "vehicleType:code"
}

func NewRepository() *Repository {
	return &Repository{
		brands: make(map[string]*domainfipe.BrandDocument),
	}
}

func key(vt domainfipe.VehicleType, code string) string {
	return fmt.Sprintf("%s:%s", vt, code)
}

func (r *Repository) GetBrands(ctx context.Context, vehicleType domainfipe.VehicleType) ([]domainfipe.Brand, time.Time, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domainfipe.Brand
	var latestSync time.Time

	for _, doc := range r.brands {
		if doc.VehicleType == string(vehicleType) {
			result = append(result, domainfipe.Brand{
				Code:        doc.Code,
				Name:        doc.Name,
				VehicleType: doc.VehicleType,
			})
			if doc.UpdatedAt.After(latestSync) {
				latestSync = doc.UpdatedAt
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, latestSync, nil
}

func (r *Repository) UpsertBrands(ctx context.Context, vehicleType domainfipe.VehicleType, brands []domainfipe.Brand, syncTime time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, b := range brands {
		k := key(vehicleType, b.Code)
		if existing, ok := r.brands[k]; ok {
			existing.Name = b.Name
			existing.UpdatedAt = syncTime
		} else {
			r.brands[k] = &domainfipe.BrandDocument{
				Code:        b.Code,
				Name:        b.Name,
				VehicleType: string(vehicleType),
				Models:      []domainfipe.ModelDocument{},
				CreatedAt:   syncTime,
				UpdatedAt:   syncTime,
			}
		}
	}
	return nil
}

func (r *Repository) GetBrandWithModels(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string) (*domainfipe.BrandDocument, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	k := key(vehicleType, brandCode)
	doc, ok := r.brands[k]
	if !ok {
		return nil, domainfipe.ErrBrandNotFound
	}

	// Deep copy to prevent external mutation
	copied := *doc
	copied.Models = make([]domainfipe.ModelDocument, len(doc.Models))
	for i, m := range doc.Models {
		copied.Models[i] = m
		copied.Models[i].Years = make([]domainfipe.YearDocument, len(m.Years))
		copy(copied.Models[i].Years, m.Years)
	}

	return &copied, nil
}

func (r *Repository) UpdateModels(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string, models []domainfipe.Model, syncTime time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	k := key(vehicleType, brandCode)
	doc, ok := r.brands[k]
	if !ok {
		return domainfipe.ErrBrandNotFound
	}

	existingMap := make(map[string]domainfipe.ModelDocument)
	for _, m := range doc.Models {
		existingMap[m.Code] = m
	}

	newModels := make([]domainfipe.ModelDocument, 0, len(models))
	for _, m := range models {
		mDoc := domainfipe.ModelDocument{
			Code: m.Code,
			Name: m.Name,
		}
		if prev, exists := existingMap[m.Code]; exists {
			mDoc.Years = prev.Years
			mDoc.YearsLastSyncAt = prev.YearsLastSyncAt
		} else {
			mDoc.Years = []domainfipe.YearDocument{}
		}
		newModels = append(newModels, mDoc)
	}

	doc.Models = newModels
	doc.ModelsLastSyncAt = syncTime
	doc.UpdatedAt = syncTime
	return nil
}

func (r *Repository) UpdateModelYears(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string, modelCode string, years []domainfipe.Year, syncTime time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	k := key(vehicleType, brandCode)
	doc, ok := r.brands[k]
	if !ok {
		return domainfipe.ErrBrandNotFound
	}

	for i := range doc.Models {
		if doc.Models[i].Code == modelCode {
			existingYearsMap := make(map[string]domainfipe.YearDocument)
			for _, y := range doc.Models[i].Years {
				existingYearsMap[y.Code] = y
			}

			newYears := make([]domainfipe.YearDocument, 0, len(years))
			for _, y := range years {
				yDoc := domainfipe.YearDocument{
					Code: y.Code,
					Name: y.Name,
				}
				if prev, exists := existingYearsMap[y.Code]; exists {
					yDoc.Price = prev.Price
					yDoc.PriceValue = prev.PriceValue
					yDoc.FIPECode = prev.FIPECode
					yDoc.Fuel = prev.Fuel
					yDoc.ReferenceMonth = prev.ReferenceMonth
					yDoc.DetailsLastSyncAt = prev.DetailsLastSyncAt
				}
				newYears = append(newYears, yDoc)
			}

			doc.Models[i].Years = newYears
			doc.Models[i].YearsLastSyncAt = syncTime
			doc.UpdatedAt = syncTime
			return nil
		}
	}

	return domainfipe.ErrModelNotFound
}

func (r *Repository) UpdateYearDetail(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string, modelCode string, yearCode string, detail domainfipe.VehicleDetail, syncTime time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	k := key(vehicleType, brandCode)
	doc, ok := r.brands[k]
	if !ok {
		return domainfipe.ErrBrandNotFound
	}

	for i := range doc.Models {
		if doc.Models[i].Code == modelCode {
			for j := range doc.Models[i].Years {
				if doc.Models[i].Years[j].Code == yearCode {
					doc.Models[i].Years[j].Price = detail.Price
					doc.Models[i].Years[j].PriceValue = detail.PriceValue
					doc.Models[i].Years[j].FIPECode = detail.CodeFipe
					doc.Models[i].Years[j].Fuel = detail.Fuel
					doc.Models[i].Years[j].ReferenceMonth = detail.ReferenceMonth
					doc.Models[i].Years[j].DetailsLastSyncAt = syncTime
					doc.UpdatedAt = syncTime
					return nil
				}
			}
			return domainfipe.ErrYearNotFound
		}
	}

	return domainfipe.ErrModelNotFound
}
