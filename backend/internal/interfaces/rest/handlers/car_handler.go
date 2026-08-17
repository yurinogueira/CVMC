package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	carusecase "cvmc/internal/application/usecase/car"
	"cvmc/internal/shared/httpx"
)

type CarHandler struct {
	service *carusecase.Service
}

func NewCarHandler(service *carusecase.Service) *CarHandler {
	return &CarHandler{service: service}
}

func (h *CarHandler) Create(w http.ResponseWriter, r *http.Request) {
	actorID := extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var input struct {
		Name            string `json:"name"`
		Manufacturer    string `json:"manufacturer"`
		Model           string `json:"model"`
		YearManufacture int    `json:"yearManufacture"`
		YearModel       int    `json:"yearModel"`
		LastMileage     int    `json:"lastMileage"`
	}
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
	})
	if err != nil {
		handleCarError(w, err)
		return
	}
	httpx.Created(w, car)
}

func (h *CarHandler) List(w http.ResponseWriter, r *http.Request) {
	actorID := extractUserID(r)
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

func (h *CarHandler) Get(w http.ResponseWriter, r *http.Request) {
	actorID := extractUserID(r)
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

func (h *CarHandler) Update(w http.ResponseWriter, r *http.Request) {
	actorID := extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	carID := r.PathValue("id")
	var input struct {
		Name            string `json:"name"`
		Manufacturer    string `json:"manufacturer"`
		Model           string `json:"model"`
		YearManufacture int    `json:"yearManufacture"`
		YearModel       int    `json:"yearModel"`
		LastMileage     int    `json:"lastMileage"`
	}
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
	})
	if err != nil {
		handleCarError(w, err)
		return
	}
	httpx.Success(w, car)
}

func (h *CarHandler) Delete(w http.ResponseWriter, r *http.Request) {
	actorID := extractUserID(r)
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

func (h *CarHandler) Share(w http.ResponseWriter, r *http.Request) {
	actorID := extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	carID := r.PathValue("id")
	var input struct {
		Email string `json:"email"`
	}
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

func (h *CarHandler) Unshare(w http.ResponseWriter, r *http.Request) {
	actorID := extractUserID(r)
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

func extractUserID(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}
