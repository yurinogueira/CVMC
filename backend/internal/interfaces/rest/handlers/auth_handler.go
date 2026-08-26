package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	portauth "cvmc/internal/application/ports/auth"
	emailport "cvmc/internal/application/ports/email"
	userport "cvmc/internal/application/ports/user"
	authusecase "cvmc/internal/application/usecase/auth"
	"cvmc/internal/shared/httpx"
)

type AuthHandler struct {
	service      *authusecase.Service
	cookieDomain string
	cookieSecure bool
}

func NewAuthHandler(users userport.Repository, hasher portauth.PasswordHasher, tokens portauth.TokenService, emailSender emailport.Sender, cookieDomain string, cookieSecure bool) *AuthHandler {
	return &AuthHandler{
		service:      authusecase.NewService(users, hasher, tokens, emailSender),
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

type ForgotPasswordRequest struct {
	Email string `json:"email" example:"yuri@cvmc.com"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" example:"a1b2c3d4e5f6..."`
	NewPassword string `json:"newPassword" example:"novaSenha12345"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" example:"a1b2c3d4e5f6..."`
}

// Register godoc
// @Summary      Cadastro de novo usuário
// @Description  Cria uma nova conta de usuário, envia e-mail de confirmação, seta cookies httpOnly com tokens e retorna dados do usuário
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
		"user": map[string]any{
			"id":            output.User.ID,
			"name":          output.User.Name,
			"email":         output.User.Email,
			"emailVerified": output.User.EmailVerified,
			"maxVehicles":   output.User.MaxVehicles,
			"createdAt":     output.User.CreatedAt,
		},
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
		"user": map[string]any{
			"id":              output.User.ID,
			"name":            output.User.Name,
			"email":           output.User.Email,
			"emailVerified":   output.User.EmailVerified,
			"emailVerifiedAt": output.User.EmailVerifiedAt,
			"maxVehicles":     output.User.MaxVehicles,
			"createdAt":       output.User.CreatedAt,
		},
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
		"user": map[string]any{
			"id":              output.User.ID,
			"name":            output.User.Name,
			"email":           output.User.Email,
			"emailVerified":   output.User.EmailVerified,
			"emailVerifiedAt": output.User.EmailVerifiedAt,
			"maxVehicles":     output.User.MaxVehicles,
			"createdAt":       output.User.CreatedAt,
		},
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
	httpx.Success(w, map[string]any{
		"id":              user.ID,
		"name":            user.Name,
		"email":           user.Email,
		"emailVerified":   user.EmailVerified,
		"emailVerifiedAt": user.EmailVerifiedAt,
		"maxVehicles":     user.MaxVehicles,
		"createdAt":       user.CreatedAt,
	})
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

// ForgotPassword godoc
// @Summary      Solicitar recuperação de senha
// @Description  Envia e-mail com instruções e link para redefinir senha sem enumeração de usuários
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload body ForgotPasswordRequest true "E-mail cadastrado"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Router       /api/v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var input ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleAuthError(w, err)
		return
	}
	if err := h.service.ForgotPassword(r.Context(), input.Email); err != nil {
		handleAuthError(w, err)
		return
	}
	httpx.Success(w, map[string]string{
		"message": "Se o e-mail informado estiver cadastrado, você receberá instruções para redefinir sua senha.",
	})
}

// ResetPassword godoc
// @Summary      Redefinir senha
// @Description  Redefine a senha do usuário utilizando o token recebido por e-mail
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload body ResetPasswordRequest true "Token e nova senha"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleAuthError(w, err)
		return
	}
	if err := h.service.ResetPassword(r.Context(), input.Token, input.NewPassword); err != nil {
		handleAuthError(w, err)
		return
	}
	httpx.Success(w, map[string]string{
		"message": "Senha redefinida com sucesso.",
	})
}

// VerifyEmail godoc
// @Summary      Validar endereço de e-mail
// @Description  Valida o e-mail do usuário utilizando o token de uso único recebido
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload body VerifyEmailRequest true "Token de validação"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var input VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleAuthError(w, err)
		return
	}
	if err := h.service.VerifyEmail(r.Context(), input.Token); err != nil {
		handleAuthError(w, err)
		return
	}
	httpx.Success(w, map[string]string{
		"message": "E-mail verificado com sucesso.",
	})
}

// ResendVerification godoc
// @Summary      Reenviar e-mail de confirmação
// @Description  Reenvia o e-mail de confirmação para o usuário autenticado
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/auth/resend-verification [post]
func (h *AuthHandler) ResendVerification(w http.ResponseWriter, r *http.Request) {
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
	if err := h.service.ResendVerification(r.Context(), user.ID); err != nil {
		handleAuthError(w, err)
		return
	}
	httpx.Success(w, map[string]string{
		"message": "E-mail de confirmação reenviado com sucesso.",
	})
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
		MaxAge:   24 * 3600, // 24 hours
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
	case errors.Is(err, authusecase.ErrAlreadyVerified):
		httpx.Error(w, http.StatusBadRequest, "Email already verified", nil)
	case errors.Is(err, authusecase.ErrTokenExpired):
		httpx.Error(w, http.StatusUnauthorized, "Token expired", nil)
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
