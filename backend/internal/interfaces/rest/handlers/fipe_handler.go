package handlers

import (
	"errors"
	"net/http"
	"strings"

	portauth "cvmc/internal/application/ports/auth"
	fipeusecase "cvmc/internal/application/usecase/fipe"
	domainfipe "cvmc/internal/domain/fipe"
	"cvmc/internal/shared/httpx"
)

type FipeHandler struct {
	service *fipeusecase.Service
	tokens  portauth.TokenService
}

func NewFipeHandler(service *fipeusecase.Service, tokens ...portauth.TokenService) *FipeHandler {
	var tokenService portauth.TokenService
	if len(tokens) > 0 {
		tokenService = tokens[0]
	}
	return &FipeHandler{service: service, tokens: tokenService}
}

func (h *FipeHandler) extractUserID(r *http.Request) string {
	if cookie, err := r.Cookie("cvmc_access_token"); err == nil && cookie.Value != "" {
		if h.tokens != nil {
			claims, err := h.tokens.ParseAccessToken(cookie.Value)
			if err == nil && claims.UserID != "" {
				return claims.UserID
			}
		}
	}
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

// ListBrands godoc
// @Summary      Listar marcas de veículos (Fipe)
// @Description  Retorna a lista de marcas para um tipo de veículo com cache multinível no MongoDB
// @Tags         Fipe
// @Produce      json
// @Security     BearerAuth
// @Param        vehicleType path string true "Tipo de veículo (cars, motorcycles, trucks)"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      503 {object} httpx.ErrorEnvelope
// @Router       /api/v1/fipe/{vehicleType}/brands [get]
func (h *FipeHandler) ListBrands(w http.ResponseWriter, r *http.Request) {
	if h.extractUserID(r) == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	vt, err := domainfipe.ParseVehicleType(r.PathValue("vehicleType"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	brands, err := h.service.GetBrands(r.Context(), vt)
	if err != nil {
		handleFipeError(w, err)
		return
	}

	httpx.Success(w, brands)
}

// ListModels godoc
// @Summary      Listar modelos de uma marca (Fipe)
// @Description  Retorna a lista de modelos de uma determinada marca com cache multinível no MongoDB
// @Tags         Fipe
// @Produce      json
// @Security     BearerAuth
// @Param        vehicleType path string true "Tipo de veículo (cars, motorcycles, trucks)"
// @Param        brandId path string true "Código da marca Fipe"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      503 {object} httpx.ErrorEnvelope
// @Router       /api/v1/fipe/{vehicleType}/brands/{brandId}/models [get]
func (h *FipeHandler) ListModels(w http.ResponseWriter, r *http.Request) {
	if h.extractUserID(r) == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	vt, err := domainfipe.ParseVehicleType(r.PathValue("vehicleType"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	brandID := r.PathValue("brandId")
	if brandID == "" {
		httpx.Error(w, http.StatusBadRequest, "brandId is required", nil)
		return
	}

	models, err := h.service.GetModels(r.Context(), vt, brandID)
	if err != nil {
		handleFipeError(w, err)
		return
	}

	httpx.Success(w, models)
}

// ListYears godoc
// @Summary      Listar anos e versões de um modelo (Fipe)
// @Description  Retorna a lista de anos/versões de um determinado modelo com cache no MongoDB
// @Tags         Fipe
// @Produce      json
// @Security     BearerAuth
// @Param        vehicleType path string true "Tipo de veículo (cars, motorcycles, trucks)"
// @Param        brandId path string true "Código da marca Fipe"
// @Param        modelId path string true "Código do modelo Fipe"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      503 {object} httpx.ErrorEnvelope
// @Router       /api/v1/fipe/{vehicleType}/brands/{brandId}/models/{modelId}/years [get]
func (h *FipeHandler) ListYears(w http.ResponseWriter, r *http.Request) {
	if h.extractUserID(r) == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	vt, err := domainfipe.ParseVehicleType(r.PathValue("vehicleType"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	brandID := r.PathValue("brandId")
	modelID := r.PathValue("modelId")
	if brandID == "" || modelID == "" {
		httpx.Error(w, http.StatusBadRequest, "brandId and modelId are required", nil)
		return
	}

	years, err := h.service.GetYears(r.Context(), vt, brandID, modelID)
	if err != nil {
		handleFipeError(w, err)
		return
	}

	httpx.Success(w, years)
}

// GetVehicleDetail godoc
// @Summary      Obter detalhes e preço Fipe de um veículo
// @Description  Retorna os detalhes completos, código Fipe, combustível e preço de mercado atualizado
// @Tags         Fipe
// @Produce      json
// @Security     BearerAuth
// @Param        vehicleType path string true "Tipo de veículo (cars, motorcycles, trucks)"
// @Param        brandId path string true "Código da marca Fipe"
// @Param        modelId path string true "Código do modelo Fipe"
// @Param        yearId path string true "Código do ano Fipe (ex: 2023-1)"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      503 {object} httpx.ErrorEnvelope
// @Router       /api/v1/fipe/{vehicleType}/brands/{brandId}/models/{modelId}/years/{yearId} [get]
func (h *FipeHandler) GetVehicleDetail(w http.ResponseWriter, r *http.Request) {
	if h.extractUserID(r) == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	vt, err := domainfipe.ParseVehicleType(r.PathValue("vehicleType"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	brandID := r.PathValue("brandId")
	modelID := r.PathValue("modelId")
	yearID := r.PathValue("yearId")
	if brandID == "" || modelID == "" || yearID == "" {
		httpx.Error(w, http.StatusBadRequest, "brandId, modelId, and yearId are required", nil)
		return
	}

	detail, err := h.service.GetVehicleDetail(r.Context(), vt, brandID, modelID, yearID)
	if err != nil {
		handleFipeError(w, err)
		return
	}

	httpx.Success(w, detail)
}

func handleFipeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domainfipe.ErrInvalidVehicleType):
		httpx.Error(w, http.StatusBadRequest, "Invalid vehicle type", nil)
	case errors.Is(err, domainfipe.ErrBrandNotFound):
		httpx.Error(w, http.StatusNotFound, "Brand not found", nil)
	case errors.Is(err, domainfipe.ErrModelNotFound):
		httpx.Error(w, http.StatusNotFound, "Model not found", nil)
	case errors.Is(err, domainfipe.ErrYearNotFound):
		httpx.Error(w, http.StatusNotFound, "Year not found", nil)
	case errors.Is(err, fipeusecase.ErrFipeUnavailable):
		httpx.Error(w, http.StatusServiceUnavailable, "Fipe service unavailable and no cache available", nil)
	default:
		httpx.Error(w, http.StatusInternalServerError, "Unexpected error retrieving Fipe data", nil)
	}
}
