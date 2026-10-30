package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	portauth "cvmc/internal/application/ports/auth"
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
				Name      string `json:"name"`
				Email     string `json:"email"`
				AvatarURL string `json:"avatarUrl"`
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

func TestUserHandlerAuthProvidersLinkUnlinkAndSetPassword(t *testing.T) {
	users := memoryuser.NewRepository()
	cars := carrepo.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("test-access-secret-32-chars-long", "test-refresh-secret-32-chars-long")
	mockVal := &mockGoogleValidator{
		payload: &portauth.GoogleTokenPayload{
			Subject:       "google-uid-carol",
			Email:         "carol@example.com",
			EmailVerified: true,
			Name:          "Carol Google",
		},
	}
	authService := auth.NewService(users, hasher, tokens, nil, mockVal)
	userHandler := NewUserHandler(users, cars, hasher, tokens, mockVal)
	ctx := context.Background()

	// Register Carol via Google (no password)
	res, err := authService.LoginWithGoogle(ctx, "carol-token")
	if err != nil {
		t.Fatalf("failed to login with google: %v", err)
	}

	// 1. GET /api/v1/user/auth-providers
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/user/auth-providers", nil)
	getReq.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: res.AccessToken})
	getRec := httptest.NewRecorder()
	userHandler.GetAuthProviders(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for GetAuthProviders, got %d: %s", getRec.Code, getRec.Body.String())
	}

	// 2. Unlink Google fails because no password
	unlinkReq := httptest.NewRequest(http.MethodPost, "/api/v1/user/auth-providers/google/unlink", nil)
	unlinkReq.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: res.AccessToken})
	unlinkRec := httptest.NewRecorder()
	userHandler.UnlinkGoogle(unlinkRec, unlinkReq)

	if unlinkRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for lockout prevention on unlink, got %d", unlinkRec.Code)
	}

	// 3. SetPassword
	setPayload := []byte(`{"newPassword":"validNewPassword123"}`)
	setReq := httptest.NewRequest(http.MethodPost, "/api/v1/user/password/set", bytes.NewReader(setPayload))
	setReq.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: res.AccessToken})
	setRec := httptest.NewRecorder()
	userHandler.SetPassword(setRec, setReq)

	if setRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for SetPassword, got %d: %s", setRec.Code, setRec.Body.String())
	}

	// 4. Now Unlink Google succeeds
	unlinkReq2 := httptest.NewRequest(http.MethodPost, "/api/v1/user/auth-providers/google/unlink", nil)
	unlinkReq2.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: res.AccessToken})
	unlinkRec2 := httptest.NewRecorder()
	userHandler.UnlinkGoogle(unlinkRec2, unlinkReq2)

	if unlinkRec2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for UnlinkGoogle after setting password, got %d: %s", unlinkRec2.Code, unlinkRec2.Body.String())
	}

	// 5. Link Google again
	linkPayload := []byte(`{"idToken":"valid-google-idtoken"}`)
	linkReq := httptest.NewRequest(http.MethodPost, "/api/v1/user/auth-providers/google/link", bytes.NewReader(linkPayload))
	linkReq.AddCookie(&http.Cookie{Name: "cvmc_access_token", Value: res.AccessToken})
	linkRec := httptest.NewRecorder()
	userHandler.LinkGoogle(linkRec, linkReq)

	if linkRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for LinkGoogle, got %d: %s", linkRec.Code, linkRec.Body.String())
	}
}
