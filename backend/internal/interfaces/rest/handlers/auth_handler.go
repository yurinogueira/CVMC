package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	portauth "cvmc/internal/application/ports/auth"
	userport "cvmc/internal/application/ports/user"
	authusecase "cvmc/internal/application/usecase/auth"
	"cvmc/internal/shared/httpx"
)

type AuthHandler struct {
	service      *authusecase.Service
	cookieDomain string
	cookieSecure bool
}

func NewAuthHandler(users userport.Repository, hasher portauth.PasswordHasher, tokens portauth.TokenService, cookieDomain string, cookieSecure bool) *AuthHandler {
	return &AuthHandler{
		service:      authusecase.NewService(users, hasher, tokens),
		cookieDomain: cookieDomain,
		cookieSecure: cookieSecure,
	}
}

type RegisterRequest struct {
	Name     string `json:"name" example:"Yuri Nogueira"`
	Email    string `json:"email" example:"yuri@cvmc.com"`
	Password string `json:"password" example:"senha12345"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"yuri@cvmc.com"`
	Password string `json:"password" example:"senha12345"`
}

// Register godoc
// @Summary      Cadastro de novo usuário
// @Description  Cria uma nova conta de usuário, seta cookies httpOnly com tokens e retorna dados do usuário
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload body RegisterRequest true "Dados de cadastro"
// @Success      201 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      409 {object} httpx.ErrorEnvelope
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleAuthError(w, err)
		return
	}
	output, err := h.service.Register(r.Context(), authusecase.RegisterInput(input))
	if err != nil {
		handleAuthError(w, err)
		return
	}
	h.setTokenCookies(w, output.AccessToken, output.RefreshToken)
	httpx.Created(w, map[string]any{
		"user": map[string]any{"id": output.User.ID, "name": output.User.Name, "email": output.User.Email, "createdAt": output.User.CreatedAt},
	})
}

// Login godoc
// @Summary      Autenticação de usuário
// @Description  Autentica com e-mail e senha, seta cookies httpOnly com tokens e retorna dados do usuário
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload body LoginRequest true "Credenciais de acesso"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleAuthError(w, err)
		return
	}
	output, err := h.service.Login(r.Context(), authusecase.LoginInput(input))
	if err != nil {
		handleAuthError(w, err)
		return
	}
	h.setTokenCookies(w, output.AccessToken, output.RefreshToken)
	httpx.Success(w, map[string]any{
		"user": map[string]any{"id": output.User.ID, "name": output.User.Name, "email": output.User.Email, "createdAt": output.User.CreatedAt},
	})
}

// Refresh godoc
// @Summary      Renovar token de acesso
// @Description  Gera um novo par de tokens a partir do cookie de refresh
// @Tags         Auth
// @Produce      json
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("cvmc_refresh_token")
	if err != nil || cookie.Value == "" {
		httpx.Error(w, http.StatusUnauthorized, "Missing refresh token", nil)
		return
	}
	output, err := h.service.Refresh(r.Context(), cookie.Value)
	if err != nil {
		handleAuthError(w, err)
		return
	}
	h.setTokenCookies(w, output.AccessToken, output.RefreshToken)
	httpx.Success(w, map[string]any{
		"user": map[string]any{"id": output.User.ID, "name": output.User.Name, "email": output.User.Email, "createdAt": output.User.CreatedAt},
	})
}

// Me godoc
// @Summary      Dados do usuário autenticado
// @Description  Retorna as informações do usuário atual com base no cookie de acesso
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/auth/me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	accessToken := h.extractAccessToken(r)
	if accessToken == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	user, err := h.service.Me(r.Context(), accessToken)
	if err != nil {
		handleAuthError(w, err)
		return
	}
	httpx.Success(w, map[string]any{"id": user.ID, "name": user.Name, "email": user.Email, "createdAt": user.CreatedAt})
}

// Logout godoc
// @Summary      Encerrar sessão
// @Description  Remove os cookies de autenticação
// @Tags         Auth
// @Produce      json
// @Success      200 {object} httpx.SuccessEnvelope
// @Router       /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_ = r
	h.clearTokenCookies(w)
	httpx.Success(w, map[string]string{"message": "logged out"})
}

func (h *AuthHandler) extractAccessToken(r *http.Request) string {
	// Try cookie first
	if cookie, err := r.Cookie("cvmc_access_token"); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	// Fallback to Authorization header for backward compatibility
	if bearer := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "); bearer != "" && bearer != r.Header.Get("Authorization") {
		return bearer
	}
	return ""
}

func (h *AuthHandler) setTokenCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "cvmc_access_token",
		Value:    accessToken,
		Path:     "/",
		Domain:   h.cookieDomain,
		MaxAge:   900, // 15 minutes
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "cvmc_refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth/refresh",
		Domain:   h.cookieDomain,
		MaxAge:   7 * 24 * 3600, // 7 days
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *AuthHandler) clearTokenCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "cvmc_access_token",
		Value:    "",
		Path:     "/",
		Domain:   h.cookieDomain,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "cvmc_refresh_token",
		Value:    "",
		Path:     "/api/v1/auth/refresh",
		Domain:   h.cookieDomain,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func handleAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, authusecase.ErrEmailInUse):
		httpx.Error(w, http.StatusConflict, "Email already in use", nil)
	case errors.Is(err, authusecase.ErrWeakPassword):
		httpx.Error(w, http.StatusBadRequest, "Password must be between 8 and 72 characters", nil)
	case errors.Is(err, authusecase.ErrInvalidInput):
		httpx.Error(w, http.StatusBadRequest, "Invalid input: please verify name, email and password", nil)
	case errors.Is(err, authusecase.ErrInvalidCredentials):
		httpx.Error(w, http.StatusUnauthorized, "Invalid credentials", nil)
	case errors.Is(err, authusecase.ErrInvalidToken):
		httpx.Error(w, http.StatusUnauthorized, "Invalid token", nil)
	case errors.Is(err, authusecase.ErrUserNotFound):
		httpx.Error(w, http.StatusNotFound, "User not found", nil)
	default:
		httpx.Error(w, http.StatusInternalServerError, "Unexpected error", nil)
	}
}
