package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	domainfipe "cvmc/internal/domain/fipe"
)

func TestClient_FetchEndpoints(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Subscription-Token")
		if token != "test-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		switch r.URL.Path {
		case "/cars/brands":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[{"code":"59","name":"VW - VolksWagen"}]`))
		case "/cars/brands/59/models":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[{"code":"5940","name":"Gol 1.0"}]`))
		case "/cars/brands/59/models/5940/years":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[{"code":"2023-1","name":"2023 Gasolina"}]`))
		case "/cars/brands/59/models/5940/years/2023-1":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"brand": "VW - VolksWagen",
				"codeFipe": "005487-9",
				"fuel": "Gasolina",
				"model": "Gol 1.0",
				"modelYear": 2023,
				"price": "R$ 55.400,00",
				"referenceMonth": "agosto de 2026",
				"vehicleType": 1
			}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "test-token")
	ctx := context.Background()

	// 1. Brands
	brands, err := client.FetchBrands(ctx, domainfipe.VehicleTypeCars)
	if err != nil {
		t.Fatalf("unexpected error fetching brands: %v", err)
	}
	if len(brands) != 1 || brands[0].Code != "59" {
		t.Fatalf("unexpected brands: %+v", brands)
	}

	// 2. Models
	models, err := client.FetchModels(ctx, domainfipe.VehicleTypeCars, "59")
	if err != nil {
		t.Fatalf("unexpected error fetching models: %v", err)
	}
	if len(models) != 1 || models[0].Code != "5940" {
		t.Fatalf("unexpected models: %+v", models)
	}

	// 3. Years
	years, err := client.FetchYears(ctx, domainfipe.VehicleTypeCars, "59", "5940")
	if err != nil {
		t.Fatalf("unexpected error fetching years: %v", err)
	}
	if len(years) != 1 || years[0].Code != "2023-1" {
		t.Fatalf("unexpected years: %+v", years)
	}

	// 4. Detail
	detail, err := client.FetchVehicleDetail(ctx, domainfipe.VehicleTypeCars, "59", "5940", "2023-1")
	if err != nil {
		t.Fatalf("unexpected error fetching detail: %v", err)
	}
	if detail.CodeFipe != "005487-9" || detail.PriceValue != 55400 {
		t.Fatalf("unexpected detail: %+v", detail)
	}
}

func TestParsePriceValue(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"R$ 55.400,00", 55400.0},
		{"R$ 1.234.567,89", 1234567.89},
		{"R$ 0,00", 0.0},
		{"invalid", 0.0},
	}

	for _, tt := range tests {
		val := parsePriceValue(tt.input)
		if val != tt.expected {
			t.Errorf("parsePriceValue(%q) = %f, want %f", tt.input, val, tt.expected)
		}
	}
}
