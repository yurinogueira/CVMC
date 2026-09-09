package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cvmc/internal/application/usecase/auth"
	carusecase "cvmc/internal/application/usecase/car"
	fuelusecase "cvmc/internal/application/usecase/fuel"
	"cvmc/internal/infrastructure/auth/bcrypt"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	carrepo "cvmc/internal/infrastructure/car/memory"
	fuelrepo "cvmc/internal/infrastructure/fuel/memory"
	memoryuser "cvmc/internal/infrastructure/user/memory"
)

func TestFuelHandlerCRUD(t *testing.T) {
	users := memoryuser.NewRepository()
	cars := carrepo.NewRepository()
	fuels := fuelrepo.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("test-access-secret-32-characters-min", "test-refresh-secret-32-characters-min")

	authService := auth.NewService(users, hasher, tokens, nil)
	carService := carusecase.NewService(cars, users)
	fuelService := fuelusecase.NewService(fuels, cars)
	handler := NewFuelHandler(fuelService, tokens)

	ctx := context.Background()

	// Register owner
	ownerReg, err := authService.Register(ctx, auth.RegisterInput{
		Name:     "Owner User",
		Email:    "owner@test.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("register owner failed: %v", err)
	}
	ownerReg.User.EmailVerified = true
	if _, err := users.Update(ctx, ownerReg.User); err != nil {
		t.Fatalf("update owner failed: %v", err)
	}

	// Register other user
	otherReg, err := authService.Register(ctx, auth.RegisterInput{
		Name:     "Other User",
		Email:    "other@test.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("register other failed: %v", err)
	}

	// Create car for owner
	carItem, err := carService.Create(ctx, ownerReg.User.ID, carusecase.CreateInput{
		Name:            "Meu Carro",
		Manufacturer:    "Volkswagen",
		Model:           "Golf",
		YearManufacture: 2020,
		YearModel:       2021,
		LastMileage:     25000,
	})
	if err != nil {
		t.Fatalf("create car failed: %v", err)
	}

	var fuelingID string

	// 1. Unauthorized create
	{
		body, _ := json.Marshal(CreateFuelingRequest{
			FuelType:      "Gasolina Aditivada",
			Liters:        45,
			PricePerLiter: 6.0,
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/cars/"+carItem.ID+"/fuelings", bytes.NewReader(body))
		req.SetPathValue("id", carItem.ID)
		rec := httptest.NewRecorder()
		handler.Create(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
		}
	}

	// 2. Success create by owner
	{
		body, _ := json.Marshal(CreateFuelingRequest{
			Date:          time.Now().UTC(),
			FuelType:      "Gasolina Aditivada",
			Liters:        45,
			PricePerLiter: 6.0,
			TotalCost:     270.0,
			IsFullTank:    true,
			GasStation:    "Posto Shell",
			Notes:         "Primeiro abastecimento completo",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/cars/"+carItem.ID+"/fuelings", bytes.NewReader(body))
		req.SetPathValue("id", carItem.ID)
		req.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: ownerReg.AccessToken})
		rec := httptest.NewRecorder()
		handler.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp struct {
			Data struct {
				ID        string  `json:"id"`
				FuelType  string  `json:"fuelType"`
				TotalCost float64 `json:"totalCost"`
			} `json:"data"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal create response: %v", err)
		}
		fuelingID = resp.Data.ID
		if fuelingID == "" {
			t.Fatalf("expected non-empty fueling ID")
		}
		if resp.Data.TotalCost != 270.0 {
			t.Fatalf("expected totalCost 270.0, got %f", resp.Data.TotalCost)
		}
	}

	// 3. List fuelings by car
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cars/"+carItem.ID+"/fuelings", nil)
		req.SetPathValue("id", carItem.ID)
		req.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: ownerReg.AccessToken})
		rec := httptest.NewRecorder()
		handler.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal list response: %v", err)
		}
		if len(resp.Data) != 1 || resp.Data[0].ID != fuelingID {
			t.Fatalf("expected 1 item with ID %s, got %+v", fuelingID, resp.Data)
		}
	}

	// 4. Get fueling by ID
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/fuelings/"+fuelingID, nil)
		req.SetPathValue("fuelingID", fuelingID)
		req.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: ownerReg.AccessToken})
		rec := httptest.NewRecorder()
		handler.Get(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}
	}

	// 5. Update by stranger (Forbidden)
	{
		body, _ := json.Marshal(UpdateFuelingRequest{
			FuelType:      "Etanol",
			Liters:        40,
			PricePerLiter: 4.0,
		})
		req := httptest.NewRequest(http.MethodPut, "/api/v1/fuelings/"+fuelingID, bytes.NewReader(body))
		req.SetPathValue("fuelingID", fuelingID)
		req.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: otherReg.AccessToken})
		rec := httptest.NewRecorder()
		handler.Update(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden, got %d", rec.Code)
		}
	}

	// 6. Update by owner (Success)
	{
		body, _ := json.Marshal(UpdateFuelingRequest{
			Date:          time.Now().UTC(),
			FuelType:      "Gasolina Comum",
			Liters:        42,
			PricePerLiter: 5.5,
			TotalCost:     231.0,
			IsFullTank:    true,
			GasStation:    "Posto Ipiranga",
		})
		req := httptest.NewRequest(http.MethodPut, "/api/v1/fuelings/"+fuelingID, bytes.NewReader(body))
		req.SetPathValue("fuelingID", fuelingID)
		req.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: ownerReg.AccessToken})
		rec := httptest.NewRecorder()
		handler.Update(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}
	}

	// 7. Delete by stranger (Forbidden)
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/fuelings/"+fuelingID, nil)
		req.SetPathValue("fuelingID", fuelingID)
		req.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: otherReg.AccessToken})
		rec := httptest.NewRecorder()
		handler.Delete(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden, got %d", rec.Code)
		}
	}

	// 8. Delete by owner (Success)
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/fuelings/"+fuelingID, nil)
		req.SetPathValue("fuelingID", fuelingID)
		req.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: ownerReg.AccessToken})
		rec := httptest.NewRecorder()
		handler.Delete(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}
	}

	// 9. Get after delete (NotFound)
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/fuelings/"+fuelingID, nil)
		req.SetPathValue("fuelingID", fuelingID)
		req.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: ownerReg.AccessToken})
		rec := httptest.NewRecorder()
		handler.Get(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected 404 NotFound, got %d", rec.Code)
		}
	}
}
