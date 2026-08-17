package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	maintenusecase "cvmc/internal/application/usecase/maintenance"
	"cvmc/internal/shared/httpx"
)

type MaintenanceHandler struct {
	service *maintenusecase.Service
}

func NewMaintenanceHandler(service *maintenusecase.Service) *MaintenanceHandler {
	return &MaintenanceHandler{service: service}
}

func (h *MaintenanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	actorID := extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	carID := r.PathValue("id")
	var input struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Date        time.Time `json:"date"`
		Mileage     int       `json:"mileage"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleMaintenanceError(w, err)
		return
	}
	maintenance, err := h.service.Create(r.Context(), actorID, carID, maintenusecase.CreateInput{
		Title:       input.Title,
		Description: input.Description,
		Date:        input.Date,
		Mileage:     input.Mileage,
	})
	if err != nil {
		handleMaintenanceError(w, err)
		return
	}
	httpx.Created(w, maintenance)
}

func (h *MaintenanceHandler) List(w http.ResponseWriter, r *http.Request) {
	actorID := extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	carID := r.PathValue("id")
	items, err := h.service.List(r.Context(), actorID, carID)
	if err != nil {
		handleMaintenanceError(w, err)
		return
	}
	httpx.Success(w, items)
}

func (h *MaintenanceHandler) Update(w http.ResponseWriter, r *http.Request) {
	actorID := extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	maintenanceID := r.PathValue("maintenanceID")
	var input struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Date        time.Time `json:"date"`
		Mileage     int       `json:"mileage"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleMaintenanceError(w, err)
		return
	}
	maintenance, err := h.service.Update(r.Context(), actorID, maintenanceID, maintenusecase.UpdateInput{
		Title:       input.Title,
		Description: input.Description,
		Date:        input.Date,
		Mileage:     input.Mileage,
	})
	if err != nil {
		handleMaintenanceError(w, err)
		return
	}
	httpx.Success(w, maintenance)
}

func (h *MaintenanceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	actorID := extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	maintenanceID := r.PathValue("maintenanceID")
	if err := h.service.Delete(r.Context(), actorID, maintenanceID); err != nil {
		handleMaintenanceError(w, err)
		return
	}
	httpx.Success(w, map[string]string{"deleted": maintenanceID})
}

func handleMaintenanceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, maintenusecase.ErrMaintenanceNotFound):
		httpx.Error(w, http.StatusNotFound, "Maintenance not found", nil)
	case errors.Is(err, maintenusecase.ErrMaintenanceForbidden):
		httpx.Error(w, http.StatusForbidden, "Forbidden", nil)
	case errors.Is(err, maintenusecase.ErrMaintenanceInvalid):
		httpx.Error(w, http.StatusBadRequest, "Invalid payload", nil)
	default:
		httpx.Error(w, http.StatusInternalServerError, "Unexpected error", nil)
	}
}
