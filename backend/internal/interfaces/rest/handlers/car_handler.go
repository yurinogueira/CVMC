package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	portauth "cvmc/internal/application/ports/auth"
	carusecase "cvmc/internal/application/usecase/car"
	"cvmc/internal/shared/httpx"
)

type CarHandler struct {
	service *carusecase.Service
	tokens  portauth.TokenService
}

func NewCarHandler(service *carusecase.Service, tokens ...portauth.TokenService) *CarHandler {
	var tokenService portauth.TokenService
	if len(tokens) > 0 {
		tokenService = tokens[0]
	}
	return &CarHandler{service: service, tokens: tokenService}
}

type CreateCarRequest struct {
	Name            string `json:"name" example:"Meu Civic"`
	Manufacturer    string `json:"manufacturer" example:"Honda"`
	Model           string `json:"model" example:"Civic Touring"`
	YearManufacture int    `json:"yearManufacture" example:"2023"`
	YearModel       int    `json:"yearModel" example:"2024"`
	LastMileage     int    `json:"lastMileage" example:"32000"`
	VehicleType     string `json:"vehicleType,omitempty" example:"cars"`
	FIPECode        string `json:"fipeCode,omitempty" example:"005487-9"`
	FIPEPrice       string `json:"fipePrice,omitempty" example:"R$ 150.000,00"`
	Fuel            string `json:"fuel,omitempty" example:"Gasolina"`
}

type UpdateCarRequest struct {
	Name            string `json:"name" example:"Meu Civic"`
	Manufacturer    string `json:"manufacturer" example:"Honda"`
	Model           string `json:"model" example:"Civic Touring"`
	YearManufacture int    `json:"yearManufacture" example:"2023"`
	YearModel       int    `json:"yearModel" example:"2024"`
	LastMileage     int    `json:"lastMileage" example:"35000"`
	VehicleType     string `json:"vehicleType,omitempty" example:"cars"`
	FIPECode        string `json:"fipeCode,omitempty" example:"005487-9"`
	FIPEPrice       string `json:"fipePrice,omitempty" example:"R$ 150.000,00"`
	Fuel            string `json:"fuel,omitempty" example:"Gasolina"`
}

type ShareCarRequest struct {
	Email string `json:"email" example:"co-motorista@cvmc.com"`
}

func (h *CarHandler) extractUserID(r *http.Request) string {
	// Try cookie first
	if cookie, err := r.Cookie("cvmc_access_token"); err == nil && cookie.Value != "" {
		if h.tokens != nil {
			claims, err := h.tokens.ParseAccessToken(cookie.Value)
			if err == nil && claims.UserID != "" {
				return claims.UserID
			}
		}
	}
	// Fallback to Authorization header
	raw := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
	if raw == "" {
		return ""
	}
	if h.tokens != nil {
		claims, err := h.tokens.ParseAccessToken(raw)
		if err == nil && claims.UserID != "" {
			return claims.UserID
		}
	}
	return ""
}

// Create godoc
// @Summary      Cadastrar novo veículo
// @Description  Registra um novo veículo associado à conta do usuário autenticado
// @Tags         Cars
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body CreateCarRequest true "Dados do veículo"
// @Success      201 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/cars [post]
func (h *CarHandler) Create(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var input CreateCarRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleCarError(w, err)
		return
	}
	car, err := h.service.Create(r.Context(), actorID, carusecase.CreateInput{
		Name:            input.Name,
		Manufacturer:    input.Manufacturer,
		Model:           input.Model,
		YearManufacture: input.YearManufacture,
		YearModel:       input.YearModel,
		LastMileage:     input.LastMileage,
		VehicleType:     input.VehicleType,
		FIPECode:        input.FIPECode,
		FIPEPrice:       input.FIPEPrice,
		Fuel:            input.Fuel,
	})
	if err != nil {
		handleCarError(w, err)
		return
	}
	httpx.Created(w, car)
}

// List godoc
// @Summary      Listar veículos
// @Description  Retorna todos os veículos pertencentes ou compartilhados com o usuário autenticado
// @Tags         Cars
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/cars [get]
func (h *CarHandler) List(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	cars, err := h.service.List(r.Context(), actorID)
	if err != nil {
		handleCarError(w, err)
		return
	}
	httpx.Success(w, cars)
}

// Get godoc
// @Summary      Obter detalhes do veículo
// @Description  Retorna os dados completos de um veículo específico pelo ID
// @Tags         Cars
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID do veículo"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/cars/{id} [get]
func (h *CarHandler) Get(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	carID := r.PathValue("id")
	car, err := h.service.Get(r.Context(), actorID, carID)
	if err != nil {
		handleCarError(w, err)
		return
	}
	httpx.Success(w, car)
}

// Update godoc
// @Summary      Atualizar veículo
// @Description  Atualiza os dados cadastrais e quilometragem do veículo
// @Tags         Cars
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID do veículo"
// @Param        payload body UpdateCarRequest true "Dados atualizados"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/cars/{id} [put]
func (h *CarHandler) Update(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	carID := r.PathValue("id")
	var input UpdateCarRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleCarError(w, err)
		return
	}
	car, err := h.service.Update(r.Context(), actorID, carID, carusecase.UpdateInput{
		Name:            input.Name,
		Manufacturer:    input.Manufacturer,
		Model:           input.Model,
		YearManufacture: input.YearManufacture,
		YearModel:       input.YearModel,
		LastMileage:     input.LastMileage,
		VehicleType:     input.VehicleType,
		FIPECode:        input.FIPECode,
		FIPEPrice:       input.FIPEPrice,
		Fuel:            input.Fuel,
	})
	if err != nil {
		handleCarError(w, err)
		return
	}
	httpx.Success(w, car)
}

// Delete godoc
// @Summary      Excluir veículo
// @Description  Remove o veículo do sistema
// @Tags         Cars
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID do veículo"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      403 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/cars/{id} [delete]
func (h *CarHandler) Delete(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	carID := r.PathValue("id")
	if err := h.service.Delete(r.Context(), actorID, carID); err != nil {
		handleCarError(w, err)
		return
	}
	httpx.Success(w, map[string]string{"deleted": carID})
}

// Share godoc
// @Summary      Compartilhar veículo
// @Description  Concede permissão de acesso ao veículo para outro usuário através do e-mail
// @Tags         Cars
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID do veículo"
// @Param        payload body ShareCarRequest true "E-mail do usuário"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/cars/{id}/share [post]
func (h *CarHandler) Share(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	carID := r.PathValue("id")
	var input ShareCarRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleCarError(w, err)
		return
	}
	car, err := h.service.Share(r.Context(), actorID, carID, input.Email)
	if err != nil {
		handleCarError(w, err)
		return
	}
	httpx.Success(w, car)
}

// Unshare godoc
// @Summary      Remover compartilhamento
// @Description  Revoga o acesso de um usuário ao veículo
// @Tags         Cars
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID do veículo"
// @Param        userID path string true "ID do usuário a desvincular"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/cars/{id}/share/{userID} [delete]
func (h *CarHandler) Unshare(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	carID := r.PathValue("id")
	userID := r.PathValue("userID")
	car, err := h.service.Unshare(r.Context(), actorID, carID, userID)
	if err != nil {
		handleCarError(w, err)
		return
	}
	httpx.Success(w, car)
}

func handleCarError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, carusecase.ErrCarNotFound):
		httpx.Error(w, http.StatusNotFound, "Car not found", nil)
	case errors.Is(err, carusecase.ErrForbidden):
		httpx.Error(w, http.StatusForbidden, "Forbidden", nil)
	case errors.Is(err, carusecase.ErrInvalidPayload):
		httpx.Error(w, http.StatusBadRequest, "Invalid payload", nil)
	case errors.Is(err, carusecase.ErrShareNotFound):
		httpx.Error(w, http.StatusNotFound, "User not found", nil)
	default:
		httpx.Error(w, http.StatusInternalServerError, "Unexpected error", nil)
	}
}
