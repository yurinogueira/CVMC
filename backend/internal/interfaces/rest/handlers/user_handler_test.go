package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cvmc/internal/application/usecase/auth"
	"cvmc/internal/infrastructure/auth/bcrypt"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	carrepo "cvmc/internal/infrastructure/car/memory"
	memoryuser "cvmc/internal/infrastructure/user/memory"
)

func TestUserHandlerGetAndUpdateProfile(t *testing.T) {
	users := memoryuser.NewRepository()
	cars := carrepo.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("test-access-secret", "test-refresh-secret")
	authService := auth.NewService(users, hasher, tokens, nil)
	userHandler := NewUserHandler(users, cars, hasher, tokens)

	registered, err := authService.Register(context.Background(), auth.RegisterInput{
		Name:     "Carolina",
		Email:    "carolina@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// 1. Get Profile
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/profile", nil)
	req.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: registered.AccessToken})
	rec := httptest.NewRecorder()
	userHandler.GetProfile(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for GetProfile, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	var getResp struct {
		Data struct {
			User struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"user"`
			VehiclesCount int `json:"vehiclesCount"`
			MaxVehicles   int `json:"maxVehicles"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if getResp.Data.User.Name != "Carolina" || getResp.Data.MaxVehicles != 3 {
		t.Fatalf("unexpected profile response: %+v", getResp)
	}

	// 2. Update Profile Name
	updatePayload := []byte(`{"name":"Carolina Souza"}`)
	upReq := httptest.NewRequest(http.MethodPut, "/api/v1/user/profile", bytes.NewReader(updatePayload))
	upReq.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: registered.AccessToken})
	upRec := httptest.NewRecorder()
	userHandler.UpdateProfile(upRec, upReq)

	if upRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for UpdateProfile, got %d. Body: %s", upRec.Code, upRec.Body.String())
	}

	// 3. Update Password
	pwdPayload := []byte(`{"currentPassword":"password123","newPassword":"newPassword123"}`)
	pwdReq := httptest.NewRequest(http.MethodPut, "/api/v1/user/password", bytes.NewReader(pwdPayload))
	pwdReq.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: registered.AccessToken})
	pwdRec := httptest.NewRecorder()
	userHandler.UpdatePassword(pwdRec, pwdReq)

	if pwdRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for UpdatePassword, got %d. Body: %s", pwdRec.Code, pwdRec.Body.String())
	}
}
