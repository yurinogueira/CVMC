package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	portauth "cvmc/internal/application/ports/auth"
	carport "cvmc/internal/application/ports/car"
	fuelport "cvmc/internal/application/ports/fuel"
	maintport "cvmc/internal/application/ports/maintenance"
	userport "cvmc/internal/application/ports/user"
	userusecase "cvmc/internal/application/usecase/user"
	"cvmc/internal/shared/httpx"
)

type UserHandler struct {
	service      *userusecase.Service
	tokens       portauth.TokenService
	cookieDomain string
	cookieSecure bool
}

func NewUserHandler(users userport.Repository, cars carport.Repository, hasher portauth.PasswordHasher, tokens portauth.TokenService, googleValidator ...portauth.GoogleTokenValidator) *UserHandler {
	return &UserHandler{
		service: userusecase.NewService(users, cars, hasher, googleValidator...),
		tokens:  tokens,
	}
}

func (h *UserHandler) SetCookieConfig(domain string, secure bool) {
	h.cookieDomain = domain
	h.cookieSecure = secure
}

func (h *UserHandler) SetCascadeDependencies(
	maintenances maintport.Repository,
	fuelings fuelport.Repository,
	cleanupTokens func(ctx context.Context, userID string) error,
) {
	h.service.SetCascadeDependencies(maintenances, fuelings, cleanupTokens)
}

type UpdateProfileRequest struct {
	Name string `json:"name" example:"Yuri Nogueira"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" example:"senhaAntiga123"`
	NewPassword     string `json:"newPassword" example:"senhaNova12345"`
}

type LinkGoogleRequest struct {
	IDToken string `json:"idToken" example:"eyJhbGciOiJSUzI1NiIsImtpZCI6..."`
}

type SetPasswordRequest struct {
	NewPassword string `json:"newPassword" example:"senhaNova12345"`
}

type DeleteAccountRequest struct {
	CurrentPassword string `json:"currentPassword,omitempty" example:"senhaAtual123"`
	Confirmation    string `json:"confirmation,omitempty" example:"EXCLUIR"`
}

// GetProfile godoc
// @Summary      Perfil do usuário autenticado
// @Description  Retorna os dados completos do usuário, incluindo status de verificação de e-mail e consumo de cota de veículos
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/user/profile [get]
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	profile, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		handleUserError(w, err)
		return
	}
	httpx.Success(w, map[string]any{
		"user": map[string]any{
			"id":              profile.User.ID,
			"name":            profile.User.Name,
			"email":           profile.User.Email,
			"emailVerified":   profile.User.EmailVerified,
			"emailVerifiedAt": profile.User.EmailVerifiedAt,
			"maxVehicles":     profile.User.MaxVehicles,
			"maxFiles":        profile.MaxFiles,
			"maxStorageBytes": profile.MaxStorageBytes,
			"createdAt":       profile.User.CreatedAt,
			"updatedAt":       profile.User.UpdatedAt,
		},
		"vehiclesCount":    profile.VehiclesCount,
		"maxVehicles":      profile.MaxVehicles,
		"filesCount":       profile.FilesCount,
		"maxFiles":         profile.MaxFiles,
		"storageUsedBytes": profile.StorageUsedBytes,
		"maxStorageBytes":  profile.MaxStorageBytes,
	})
}

// UpdateProfile godoc
// @Summary      Atualizar perfil do usuário
// @Description  Atualiza informações do perfil do usuário autenticado (como o nome)
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body UpdateProfileRequest true "Dados para atualização"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/user/profile [put]
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var input UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleUserError(w, err)
		return
	}
	updated, err := h.service.UpdateProfile(r.Context(), userID, input.Name)
	if err != nil {
		handleUserError(w, err)
		return
	}
	httpx.Success(w, map[string]any{
		"user": map[string]any{
			"id":              updated.ID,
			"name":            updated.Name,
			"email":           updated.Email,
			"emailVerified":   updated.EmailVerified,
			"emailVerifiedAt": updated.EmailVerifiedAt,
			"maxVehicles":     updated.MaxVehicles,
			"createdAt":       updated.CreatedAt,
			"updatedAt":       updated.UpdatedAt,
		},
	})
}

// UpdatePassword godoc
// @Summary      Alterar senha do usuário
// @Description  Altera a senha do usuário autenticado exigindo a senha atual
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body UpdatePasswordRequest true "Senha atual e nova senha"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/user/password [put]
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var input UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleUserError(w, err)
		return
	}
	if err := h.service.UpdatePassword(r.Context(), userID, input.CurrentPassword, input.NewPassword); err != nil {
		handleUserError(w, err)
		return
	}
	httpx.Success(w, map[string]string{
		"message": "Senha alterada com sucesso.",
	})
}

// GetAuthProviders godoc
// @Summary      Provedores de autenticação conectados
// @Description  Retorna o status dos provedores de login vinculados à conta do usuário (senha, Google)
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/user/auth-providers [get]
func (h *UserHandler) GetAuthProviders(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	providers, err := h.service.GetAuthProviders(r.Context(), userID)
	if err != nil {
		handleUserError(w, err)
		return
	}
	httpx.Success(w, providers)
}

// LinkGoogle godoc
// @Summary      Vincular conta Google
// @Description  Vincula uma conta Google à conta do usuário autenticado validando o ID Token
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body LinkGoogleRequest true "Google ID Token"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      409 {object} httpx.ErrorEnvelope
// @Router       /api/v1/user/auth-providers/google/link [post]
func (h *UserHandler) LinkGoogle(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var input LinkGoogleRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleUserError(w, err)
		return
	}
	if err := h.service.LinkGoogleAccount(r.Context(), userID, input.IDToken); err != nil {
		handleUserError(w, err)
		return
	}
	httpx.Success(w, map[string]string{
		"message": "Conta Google vinculada com sucesso.",
	})
}

// UnlinkGoogle godoc
// @Summary      Desvincular conta Google
// @Description  Desvincula a conta Google da conta do usuário autenticado (requer senha cadastrada para evitar lockout)
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/user/auth-providers/google/unlink [post]
func (h *UserHandler) UnlinkGoogle(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	if err := h.service.UnlinkGoogleAccount(r.Context(), userID); err != nil {
		handleUserError(w, err)
		return
	}
	httpx.Success(w, map[string]string{
		"message": "Conta Google desvinculada com sucesso.",
	})
}

// SetPassword godoc
// @Summary      Definir senha inicial para conta social
// @Description  Define a senha para um usuário que se cadastrou via Google e ainda não possui senha local
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body SetPasswordRequest true "Nova senha"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/user/password/set [post]
func (h *UserHandler) SetPassword(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var input SetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleUserError(w, err)
		return
	}
	if err := h.service.SetPassword(r.Context(), userID, input.NewPassword); err != nil {
		handleUserError(w, err)
		return
	}
	httpx.Success(w, map[string]string{
		"message": "Senha cadastrada com sucesso.",
	})
}

// DeleteAccount godoc
// @Summary      Excluir conta permanentemente
// @Description  Remove permanentemente a conta do usuário autenticado e todos os seus dados vinculados (veículos próprios, manutenções, abastecimentos e tokens de dispositivos)
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body DeleteAccountRequest true "Confirmação de exclusão (senha atual ou palavra EXCLUIR)"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/user/account [delete]
func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var input DeleteAccountRequest
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil && !errors.Is(err, io.EOF) {
			handleUserError(w, userusecase.ErrInvalidInput)
			return
		}
	}
	if err := h.service.DeleteAccount(r.Context(), userID, input.CurrentPassword, input.Confirmation); err != nil {
		handleUserError(w, err)
		return
	}
	h.clearTokenCookies(w)
	httpx.Success(w, map[string]string{
		"message": "Sua conta e todos os dados vinculados foram excluídos permanentemente.",
	})
}

func (h *UserHandler) clearTokenCookies(w http.ResponseWriter) {
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

func (h *UserHandler) extractUserID(r *http.Request) string {
	token := ""
	if cookie, err := r.Cookie("cvmc_access_token"); err == nil && cookie.Value != "" {
		token = cookie.Value
	} else if bearer := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "); bearer != "" && bearer != r.Header.Get("Authorization") {
		token = bearer
	}
	if token == "" {
		return ""
	}
	claims, err := h.tokens.ParseAccessToken(token)
	if err != nil {
		return ""
	}
	return claims.UserID
}

func handleUserError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, userusecase.ErrInvalidCurrentPassword):
		httpx.Error(w, http.StatusBadRequest, "Senha atual incorreta", nil)
	case errors.Is(err, userusecase.ErrWeakPassword):
		httpx.Error(w, http.StatusBadRequest, "A nova senha deve conter entre 8 e 72 caracteres", nil)
	case errors.Is(err, userusecase.ErrInvalidInput):
		httpx.Error(w, http.StatusBadRequest, "Dados inválidos: verifique os campos informados", nil)
	case errors.Is(err, userusecase.ErrInvalidGoogleToken):
		httpx.Error(w, http.StatusUnauthorized, "Token do Google inválido ou expirado", nil)
	case errors.Is(err, userusecase.ErrGoogleEmailMismatch):
		httpx.Error(w, http.StatusBadRequest, "O e-mail da conta Google não corresponde ao e-mail da sua conta", nil)
	case errors.Is(err, userusecase.ErrGoogleAlreadyLinked):
		httpx.Error(w, http.StatusConflict, "Esta conta Google já está vinculada a outro usuário", nil)
	case errors.Is(err, userusecase.ErrGoogleNotLinked):
		httpx.Error(w, http.StatusBadRequest, "Conta Google não está vinculada a este usuário", nil)
	case errors.Is(err, userusecase.ErrLockoutPrevented):
		httpx.Error(w, http.StatusBadRequest, "Você deve cadastrar uma senha antes de desvincular sua conta Google", nil)
	case errors.Is(err, userusecase.ErrPasswordAlreadySet):
		httpx.Error(w, http.StatusBadRequest, "Sua conta já possui uma senha. Utilize a alteração de senha", nil)
	case errors.Is(err, userusecase.ErrAccountDeletionConfirmation):
		httpx.Error(w, http.StatusBadRequest, "Digite EXCLUIR para confirmar a exclusão permanente da conta", nil)
	case errors.Is(err, userusecase.ErrUserNotFound):
		httpx.Error(w, http.StatusNotFound, "Usuário não encontrado", nil)
	default:
		httpx.Error(w, http.StatusInternalServerError, "Erro inesperado", nil)
	}
}
