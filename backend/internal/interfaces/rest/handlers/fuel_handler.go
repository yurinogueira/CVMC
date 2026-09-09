package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	portauth "cvmc/internal/application/ports/auth"
	fuelusecase "cvmc/internal/application/usecase/fuel"
	"cvmc/internal/shared/httpx"
)

type FuelHandler struct {
	service *fuelusecase.Service
	tokens  portauth.TokenService
}

func NewFuelHandler(service *fuelusecase.Service, tokens ...portauth.TokenService) *FuelHandler {
	var tokenService portauth.TokenService
	if len(tokens) > 0 {
		tokenService = tokens[0]
	}
	return &FuelHandler{service: service, tokens: tokenService}
}

type CreateFuelingRequest struct {
	Date          time.Time `json:"date" example:"2026-09-08T15:04:05Z"`
	FuelType      string    `json:"fuelType" example:"Gasolina Comum"`
	Liters        float64   `json:"liters" example:"45.5"`
	PricePerLiter float64   `json:"pricePerLiter" example:"5.89"`
	TotalCost     float64   `json:"totalCost" example:"267.995"`
	IsFullTank    bool      `json:"isFullTank" example:"true"`
	GasStation    string    `json:"gasStation,omitempty" example:"Posto Shell Central"`
	Notes         string    `json:"notes,omitempty" example:"Tanque cheio após viagem"`
}

type UpdateFuelingRequest struct {
	Date          time.Time `json:"date" example:"2026-09-08T15:04:05Z"`
	FuelType      string    `json:"fuelType" example:"Gasolina Comum"`
	Liters        float64   `json:"liters" example:"45.5"`
	PricePerLiter float64   `json:"pricePerLiter" example:"5.89"`
	TotalCost     float64   `json:"totalCost" example:"267.995"`
	IsFullTank    bool      `json:"isFullTank" example:"true"`
	GasStation    string    `json:"gasStation,omitempty" example:"Posto Shell Central"`
	Notes         string    `json:"notes,omitempty" example:"Tanque cheio após viagem"`
}

func (h *FuelHandler) extractUserID(r *http.Request) string {
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
// @Summary      Cadastrar abastecimento
// @Description  Registra um novo abastecimento de combustível para o veículo especificado
// @Tags         Fuelings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID do veículo"
// @Param        payload body CreateFuelingRequest true "Dados do abastecimento"
// @Success      201 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      403 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/cars/{id}/fuelings [post]
func (h *FuelHandler) Create(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	carID := r.PathValue("id")
	var input CreateFuelingRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleFuelError(w, err)
		return
	}

	fueling, err := h.service.Create(r.Context(), actorID, carID, fuelusecase.CreateInput{
		Date:          input.Date,
		FuelType:      input.FuelType,
		Liters:        input.Liters,
		PricePerLiter: input.PricePerLiter,
		TotalCost:     input.TotalCost,
		IsFullTank:    input.IsFullTank,
		GasStation:    input.GasStation,
		Notes:         input.Notes,
	})
	if err != nil {
		handleFuelError(w, err)
		return
	}
	httpx.Created(w, fueling)
}

// List godoc
// @Summary      Listar abastecimentos do veículo
// @Description  Retorna o histórico cronológico de abastecimentos de combustível de um veículo
// @Tags         Fuelings
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID do veículo"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      403 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/cars/{id}/fuelings [get]
func (h *FuelHandler) List(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	carID := r.PathValue("id")
	items, err := h.service.List(r.Context(), actorID, carID)
	if err != nil {
		handleFuelError(w, err)
		return
	}
	httpx.Success(w, items)
}

// Get godoc
// @Summary      Obter detalhes de um abastecimento
// @Description  Retorna os detalhes de um registro específico de abastecimento
// @Tags         Fuelings
// @Produce      json
// @Security     BearerAuth
// @Param        fuelingID path string true "ID do abastecimento"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      403 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/fuelings/{fuelingID} [get]
func (h *FuelHandler) Get(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	fuelingID := r.PathValue("fuelingID")
	item, err := h.service.Get(r.Context(), actorID, fuelingID)
	if err != nil {
		handleFuelError(w, err)
		return
	}
	httpx.Success(w, item)
}

// Update godoc
// @Summary      Atualizar abastecimento
// @Description  Atualiza os dados de um abastecimento existente
// @Tags         Fuelings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        fuelingID path string true "ID do abastecimento"
// @Param        payload body UpdateFuelingRequest true "Dados atualizados do abastecimento"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      403 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/fuelings/{fuelingID} [put]
func (h *FuelHandler) Update(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	fuelingID := r.PathValue("fuelingID")
	var input UpdateFuelingRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleFuelError(w, err)
		return
	}

	fueling, err := h.service.Update(r.Context(), actorID, fuelingID, fuelusecase.UpdateInput{
		Date:          input.Date,
		FuelType:      input.FuelType,
		Liters:        input.Liters,
		PricePerLiter: input.PricePerLiter,
		TotalCost:     input.TotalCost,
		IsFullTank:    input.IsFullTank,
		GasStation:    input.GasStation,
		Notes:         input.Notes,
	})
	if err != nil {
		handleFuelError(w, err)
		return
	}
	httpx.Success(w, fueling)
}

// Delete godoc
// @Summary      Excluir abastecimento
// @Description  Remove o registro de abastecimento
// @Tags         Fuelings
// @Produce      json
// @Security     BearerAuth
// @Param        fuelingID path string true "ID do abastecimento"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      403 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/fuelings/{fuelingID} [delete]
func (h *FuelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	fuelingID := r.PathValue("fuelingID")
	if err := h.service.Delete(r.Context(), actorID, fuelingID); err != nil {
		handleFuelError(w, err)
		return
	}
	httpx.Success(w, map[string]string{"deleted": fuelingID})
}

func handleFuelError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, fuelusecase.ErrFuelingNotFound):
		httpx.Error(w, http.StatusNotFound, "Fueling not found", nil)
	case errors.Is(err, fuelusecase.ErrFuelingForbidden):
		httpx.Error(w, http.StatusForbidden, "Forbidden", nil)
	case errors.Is(err, fuelusecase.ErrFuelingInvalid):
		httpx.Error(w, http.StatusBadRequest, "Invalid payload", nil)
	default:
		httpx.Error(w, http.StatusInternalServerError, "Unexpected error", nil)
	}
}
