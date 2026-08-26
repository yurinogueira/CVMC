package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	portauth "cvmc/internal/application/ports/auth"
	carport "cvmc/internal/application/ports/car"
	userport "cvmc/internal/application/ports/user"
	userusecase "cvmc/internal/application/usecase/user"
	"cvmc/internal/shared/httpx"
)

type UserHandler struct {
	service *userusecase.Service
	tokens  portauth.TokenService
}

func NewUserHandler(users userport.Repository, cars carport.Repository, hasher portauth.PasswordHasher, tokens portauth.TokenService) *UserHandler {
	return &UserHandler{
		service: userusecase.NewService(users, cars, hasher),
		tokens:  tokens,
	}
}

type UpdateProfileRequest struct {
	Name string `json:"name" example:"Yuri Nogueira"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" example:"senhaAntiga123"`
	NewPassword     string `json:"newPassword" example:"senhaNova12345"`
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
			"createdAt":       profile.User.CreatedAt,
			"updatedAt":       profile.User.UpdatedAt,
		},
		"vehiclesCount": profile.VehiclesCount,
		"maxVehicles":   profile.MaxVehicles,
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
	case errors.Is(err, userusecase.ErrUserNotFound):
		httpx.Error(w, http.StatusNotFound, "Usuário não encontrado", nil)
	default:
		httpx.Error(w, http.StatusInternalServerError, "Erro inesperado", nil)
	}
}
