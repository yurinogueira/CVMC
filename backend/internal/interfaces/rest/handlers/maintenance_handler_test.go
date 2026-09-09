package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cvmc/internal/application/usecase/auth"
	carusecase "cvmc/internal/application/usecase/car"
	maintenusecase "cvmc/internal/application/usecase/maintenance"
	"cvmc/internal/infrastructure/auth/bcrypt"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	carrepo "cvmc/internal/infrastructure/car/memory"
	maintrepo "cvmc/internal/infrastructure/maintenance/memory"
	memoryuser "cvmc/internal/infrastructure/user/memory"
)

func TestMaintenanceHandlerCRUD(t *testing.T) {
	users := memoryuser.NewRepository()
	cars := carrepo.NewRepository()
	maints := maintrepo.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("test-access-secret", "test-refresh-secret")

	authService := auth.NewService(users, hasher, tokens, nil)
	carService := carusecase.NewService(cars, users)
	maintService := maintenusecase.NewService(maints, cars)
	handler := NewMaintenanceHandler(maintService, tokens)

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
		Manufacturer:    "Toyota",
		Model:           "Corolla",
		YearManufacture: 2020,
		YearModel:       2021,
		LastMileage:     30000,
	})
	if err != nil {
		t.Fatalf("create car failed: %v", err)
	}

	// 1. Unauthorized when no token
	{
		req := httptest.NewRequest(http.MethodPost, "/api/v1/cars/"+carItem.ID+"/maintenances", bytes.NewReader([]byte(`{"title":"Test"}`)))
		req.SetPathValue("id", carItem.ID)
		rec := httptest.NewRecorder()
		handler.Create(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
		}
	}

	// 2. Create Maintenance
	var maintID string
	{
		createBody := []byte(`{
			"title": "Revisão dos 30.000 km",
			"description": "Troca de óleo e filtro",
			"date": "2026-08-10T00:00:00Z",
			"mileage": 30500,
			"types": ["Óleo de Motor"],
			"cost": 320.00
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/cars/"+carItem.ID+"/maintenances", bytes.NewReader(createBody))
		req.SetPathValue("id", carItem.ID)
		req.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: ownerReg.AccessToken})
		rec := httptest.NewRecorder()
		handler.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp struct {
			Data struct {
				ID    string `json:"id"`
				Title string `json:"title"`
			} `json:"data"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal create response: %v", err)
		}
		maintID = resp.Data.ID
		if maintID == "" {
			t.Fatalf("expected non-empty maintenance ID")
		}
	}

	// 3. List Maintenances
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/cars/"+carItem.ID+"/maintenances", nil)
		req.SetPathValue("id", carItem.ID)
		req.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: ownerReg.AccessToken})
		rec := httptest.NewRecorder()
		handler.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}
	}

	// 4. Get Maintenance by ID
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/maintenances/"+maintID, nil)
		req.SetPathValue("maintenanceID", maintID)
		req.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: ownerReg.AccessToken})
		rec := httptest.NewRecorder()
		handler.Get(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK for Get, got %d: %s", rec.Code, rec.Body.String())
		}

		// Forbidden when accessed by other user
		reqOther := httptest.NewRequest(http.MethodGet, "/api/v1/maintenances/"+maintID, nil)
		reqOther.SetPathValue("maintenanceID", maintID)
		reqOther.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: otherReg.AccessToken})
		recOther := httptest.NewRecorder()
		handler.Get(recOther, reqOther)

		if recOther.Code != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden for other user Get, got %d", recOther.Code)
		}
	}

	// 5. Update Maintenance
	{
		updateBody := []byte(`{
			"title": "Revisão dos 30.000 km Completa",
			"description": "Atualizado com higienização",
			"date": "2026-08-10T00:00:00Z",
			"mileage": 31000,
			"types": ["Óleo de Motor", "Filtro do Ar-Condicionado"],
			"cost": 450.00
		}`)

		// Forbidden for other user
		reqForbidden := httptest.NewRequest(http.MethodPut, "/api/v1/maintenances/"+maintID, bytes.NewReader(updateBody))
		reqForbidden.SetPathValue("maintenanceID", maintID)
		reqForbidden.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: otherReg.AccessToken})
		recForbidden := httptest.NewRecorder()
		handler.Update(recForbidden, reqForbidden)
		if recForbidden.Code != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden for update by other user, got %d", recForbidden.Code)
		}

		// Success for owner
		req := httptest.NewRequest(http.MethodPut, "/api/v1/maintenances/"+maintID, bytes.NewReader(updateBody))
		req.SetPathValue("maintenanceID", maintID)
		req.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: ownerReg.AccessToken})
		rec := httptest.NewRecorder()
		handler.Update(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK for update, got %d: %s", rec.Code, rec.Body.String())
		}
	}

	// 6. Delete Maintenance
	{
		// Forbidden for other user
		reqForbidden := httptest.NewRequest(http.MethodDelete, "/api/v1/maintenances/"+maintID, nil)
		reqForbidden.SetPathValue("maintenanceID", maintID)
		reqForbidden.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: otherReg.AccessToken})
		recForbidden := httptest.NewRecorder()
		handler.Delete(recForbidden, reqForbidden)
		if recForbidden.Code != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden for delete by other user, got %d", recForbidden.Code)
		}

		// Success for owner
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/maintenances/"+maintID, nil)
		req.SetPathValue("maintenanceID", maintID)
		req.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: ownerReg.AccessToken})
		rec := httptest.NewRecorder()
		handler.Delete(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK for delete, got %d: %s", rec.Code, rec.Body.String())
		}

		// 7. Not found when trying to delete or get deleted maintenance
		reqNotFound := httptest.NewRequest(http.MethodGet, "/api/v1/maintenances/"+maintID, nil)
		reqNotFound.SetPathValue("maintenanceID", maintID)
		reqNotFound.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: ownerReg.AccessToken})
		recNotFound := httptest.NewRecorder()
		handler.Get(recNotFound, reqNotFound)

		if recNotFound.Code != http.StatusNotFound {
			t.Fatalf("expected 404 Not Found for deleted item, got %d", recNotFound.Code)
		}
	}
}
