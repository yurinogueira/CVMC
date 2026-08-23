package mongo

import (
	"testing"
	"time"

	domaincar "cvmc/internal/domain/car"
)

func TestToCarDocAndToDomainCar(t *testing.T) {
	now := time.Now().UTC()
	original := domaincar.Car{
		ID:              "car-123",
		OwnerID:         "user-456",
		Name:            "Meu Civic",
		Manufacturer:    "Honda",
		Model:           "Civic Touring",
		YearManufacture: 2021,
		YearModel:       2022,
		LastMileage:     45000,
		SharedWith:      []string{"user-789"},
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	doc := toCarDoc(original)
	if doc.ID != original.ID || doc.OwnerID != original.OwnerID || doc.Name != original.Name {
		t.Fatalf("mismatch in toCarDoc conversion: %+v", doc)
	}

	domain := toDomainCar(doc)
	if domain.ID != original.ID || domain.Model != original.Model || len(domain.SharedWith) != 1 {
		t.Fatalf("mismatch in toDomainCar conversion: %+v", domain)
	}
}

func TestNilSharedWithHandled(t *testing.T) {
	car := domaincar.Car{
		ID:         "car-1",
		SharedWith: nil,
	}
	doc := toCarDoc(car)
	if doc.SharedWith == nil {
		t.Fatalf("expected non-nil sharedWith in doc")
	}

	domain := toDomainCar(doc)
	if domain.SharedWith == nil {
		t.Fatalf("expected non-nil sharedWith in domain")
	}
}
