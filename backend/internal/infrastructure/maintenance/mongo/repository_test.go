package mongo

import (
	"testing"
	"time"

	domainmaintenance "cvmc/internal/domain/maintenance"
)

func TestToMaintenanceDocAndToDomainMaintenance(t *testing.T) {
	now := time.Now().UTC()
	original := domainmaintenance.Maintenance{
		ID:          "maint-123",
		CarID:       "car-456",
		Title:       "Troca de óleo e filtro",
		Description: "Óleo 0W20 Sintético + Filtro original",
		Date:        now,
		Mileage:     50000,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	doc := toMaintenanceDoc(original)
	if doc.ID != original.ID || doc.CarID != original.CarID || doc.Title != original.Title || doc.Mileage != 50000 {
		t.Fatalf("mismatch in toMaintenanceDoc: %+v", doc)
	}

	domain := toDomainMaintenance(doc)
	if domain.ID != original.ID || domain.Description != original.Description || domain.Date != original.Date {
		t.Fatalf("mismatch in toDomainMaintenance: %+v", domain)
	}
}
