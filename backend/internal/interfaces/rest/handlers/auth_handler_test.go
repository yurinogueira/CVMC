package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	portauth "cvmc/internal/application/ports/auth"
	"cvmc/internal/infrastructure/auth/bcrypt"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	memoryuser "cvmc/internal/infrastructure/user/memory"
)

func setupTestAuthHandler() *AuthHandler {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("test-access-secret", "test-refresh-secret")
	return NewAuthHandler(users, hasher, tokens, nil, "", false)
}

func TestAuthHandlerRegisterLoginRefreshLogoutMe(t *testing.T) {
	handler := setupTestAuthHandler()

	// 1. Test Register
	registerPayload := []byte(`{"name":"Yuri Teste","email":"yuri@example.com","password":"password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(registerPayload))
	rec := httptest.NewRecorder()
	handler.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	cookies := rec.Result().Cookies()
	var accessTokenCookie, refreshTokenCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "cvmc_access_token" {
			accessTokenCookie = c
		}
		if c.Name == "cvmc_refresh_token" {
			refreshTokenCookie = c
		}
	}

	if accessTokenCookie == nil || accessTokenCookie.Value == "" {
		t.Fatalf("expected cvmc_access_token cookie to be set")
	}
	if accessTokenCookie.MaxAge != 24*3600 {
		t.Fatalf("expected cvmc_access_token MaxAge to be 86400 (24h), got %d", accessTokenCookie.MaxAge)
	}
	if refreshTokenCookie == nil || refreshTokenCookie.Value == "" {
		t.Fatalf("expected cvmc_refresh_token cookie to be set")
	}
	if refreshTokenCookie.MaxAge != 7*24*3600 {
		t.Fatalf("expected cvmc_refresh_token MaxAge to be 604800 (7d), got %d", refreshTokenCookie.MaxAge)
	}

	var regResp struct {
		Success bool `json:"success"`
		Data    struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
			User         struct {
				Email string `json:"email"`
			} `json:"user"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &regResp); err != nil || regResp.Data.AccessToken == "" || regResp.Data.RefreshToken == "" {
		t.Fatalf("expected accessToken and refreshToken in register response JSON, got %+v", regResp)
	}

	// 2. Test Me with access token cookie
	meReq := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	meReq.AddCookie(accessTokenCookie)
	meRec := httptest.NewRecorder()
	handler.Me(meRec, meReq)

	if meRec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK for /me, got %d. Body: %s", meRec.Code, meRec.Body.String())
	}

	var meResp struct {
		Success bool `json:"success"`
		Data    struct {
			Email string `json:"email"`
		} `json:"data"`
	}
	if err := json.Unmarshal(meRec.Body.Bytes(), &meResp); err != nil || meResp.Data.Email != "yuri@example.com" {
		t.Fatalf("expected email yuri@example.com in /me response, got %+v", meResp)
	}

	// 3. Test Refresh with refresh token cookie
	refreshReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	refreshReq.AddCookie(refreshTokenCookie)
	refreshRec := httptest.NewRecorder()
	handler.Refresh(refreshRec, refreshReq)

	if refreshRec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK for /refresh, got %d. Body: %s", refreshRec.Code, refreshRec.Body.String())
	}

	refreshedCookies := refreshRec.Result().Cookies()
	var newAccessCookie, newRefreshCookie *http.Cookie
	for _, c := range refreshedCookies {
		if c.Name == "cvmc_access_token" {
			newAccessCookie = c
		}
		if c.Name == "cvmc_refresh_token" {
			newRefreshCookie = c
		}
	}
	if newAccessCookie == nil || newAccessCookie.MaxAge != 24*3600 {
		t.Fatalf("expected new access cookie with 24h MaxAge")
	}
	if newRefreshCookie == nil || newRefreshCookie.MaxAge != 7*24*3600 {
		t.Fatalf("expected new refresh cookie with 7d MaxAge")
	}

	// 3.1 Test Refresh with Authorization Bearer header (Mobile fallback)
	refreshBearerReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	refreshBearerReq.Header.Set("Authorization", "Bearer "+newRefreshCookie.Value)
	refreshBearerRec := httptest.NewRecorder()
	handler.Refresh(refreshBearerRec, refreshBearerReq)

	if refreshBearerRec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK for /refresh with Bearer header, got %d. Body: %s", refreshBearerRec.Code, refreshBearerRec.Body.String())
	}

	// 4. Test Logout
	logoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	logoutRec := httptest.NewRecorder()
	handler.Logout(logoutRec, logoutReq)

	if logoutRec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK for /logout, got %d", logoutRec.Code)
	}

	clearedCookies := logoutRec.Result().Cookies()
	for _, c := range clearedCookies {
		if (c.Name == "cvmc_access_token" || c.Name == "cvmc_refresh_token") && c.MaxAge != -1 {
			t.Fatalf("expected cookie %s to have MaxAge -1, got %d", c.Name, c.MaxAge)
		}
	}
}

func TestAuthHandlerForgotPasswordAndVerification(t *testing.T) {
	handler := setupTestAuthHandler()

	// Forgot password for unknown email (should return 200 OK)
	fpPayload := []byte(`{"email":"unknown@example.com"}`)
	fpReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", bytes.NewReader(fpPayload))
	fpRec := httptest.NewRecorder()
	handler.ForgotPassword(fpRec, fpReq)

	if fpRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for forgot-password, got %d", fpRec.Code)
	}

	// Verify email with invalid token (should return 401)
	vePayload := []byte(`{"token":"invalid-token"}`)
	veReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/verify-email", bytes.NewReader(vePayload))
	veRec := httptest.NewRecorder()
	handler.VerifyEmail(veRec, veReq)

	if veRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for invalid verify token, got %d", veRec.Code)
	}
}

func TestAuthHandlerUnauthorizedWithoutTokens(t *testing.T) {
	handler := setupTestAuthHandler()

	// /me without token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	rec := httptest.NewRecorder()
	handler.Me(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 Unauthorized for /me without token, got %d", rec.Code)
	}

	// /refresh without cookie
	reqRefresh := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	recRefresh := httptest.NewRecorder()
	handler.Refresh(recRefresh, reqRefresh)

	if recRefresh.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 Unauthorized for /refresh without cookie, got %d", recRefresh.Code)
	}
}

type mockGoogleValidator struct {
	payload *portauth.GoogleTokenPayload
	err     error
}

func (m *mockGoogleValidator) Validate(ctx context.Context, idToken string) (*portauth.GoogleTokenPayload, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.payload, nil
}

func TestAuthHandlerGoogleAuth(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("test-access-secret-32-chars-long", "test-refresh-secret-32-chars-long")
	mockVal := &mockGoogleValidator{
		payload: &portauth.GoogleTokenPayload{
			Subject:       "google-uid-test",
			Email:         "google-handler@test.com",
			EmailVerified: true,
			Name:          "Google Handler Test",
		},
	}
	handler := NewAuthHandler(users, hasher, tokens, nil, "", false, mockVal)

	payload := []byte(`{"idToken":"valid-google-idtoken"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/google", bytes.NewReader(payload))
	rec := httptest.NewRecorder()
	handler.GoogleAuth(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK for GoogleAuth, got %d: %s", rec.Code, rec.Body.String())
	}

	cookies := rec.Result().Cookies()
	var accessTokenCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "cvmc_access_token" {
			accessTokenCookie = c
		}
	}
	if accessTokenCookie == nil || accessTokenCookie.Value == "" {
		t.Fatalf("expected access token cookie to be set on GoogleAuth")
	}

	// Invalid token returns 401
	mockVal.err = errors.New("invalid token")
	reqInvalid := httptest.NewRequest(http.MethodPost, "/api/v1/auth/google", bytes.NewReader(payload))
	recInvalid := httptest.NewRecorder()
	handler.GoogleAuth(recInvalid, reqInvalid)

	if recInvalid.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for invalid token, got %d", recInvalid.Code)
	}
}
