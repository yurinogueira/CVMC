package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	portauth "cvmc/internal/application/ports/auth"
	maintenusecase "cvmc/internal/application/usecase/maintenance"
	domainmaintenance "cvmc/internal/domain/maintenance"
	"cvmc/internal/shared/httpx"
)

type MaintenanceHandler struct {
	service *maintenusecase.Service
	tokens  portauth.TokenService
}

func NewMaintenanceHandler(service *maintenusecase.Service, tokens ...portauth.TokenService) *MaintenanceHandler {
	var tokenService portauth.TokenService
	if len(tokens) > 0 {
		tokenService = tokens[0]
	}
	return &MaintenanceHandler{service: service, tokens: tokenService}
}

type AttachmentRequest struct {
	ID        string    `json:"id" example:"att-1"`
	Name      string    `json:"name" example:"recibo_oleo.pdf"`
	Size      int64     `json:"size" example:"1048576"`
	MimeType  string    `json:"mimeType" example:"application/pdf"`
	DataUrl   string    `json:"dataUrl" example:"data:application/pdf;base64,..."`
	CreatedAt time.Time `json:"createdAt" example:"2026-09-04T12:00:00Z"`
}

type CreateMaintenanceRequest struct {
	Title       string              `json:"title" example:"Troca de Óleo e Filtro"`
	Description string              `json:"description" example:"Óleo 0W20 Sintético + Filtro original"`
	Date        time.Time           `json:"date" example:"2026-08-23T00:00:00Z"`
	Mileage     int                 `json:"mileage" example:"32000"`
	Types       []string            `json:"types,omitempty" example:"[\"Óleo de Motor\",\"Filtro do Óleo de Motor\"]"`
	Cost        *float64            `json:"cost,omitempty" example:"350.50"`
	Attachments []AttachmentRequest `json:"attachments,omitempty"`
}

type UpdateMaintenanceRequest struct {
	Title       string              `json:"title" example:"Troca de Óleo e Filtro"`
	Description string              `json:"description" example:"Óleo 0W20 Sintético + Filtro original"`
	Date        time.Time           `json:"date" example:"2026-08-23T00:00:00Z"`
	Mileage     int                 `json:"mileage" example:"32000"`
	Types       []string            `json:"types,omitempty" example:"[\"Óleo de Motor\",\"Filtro do Óleo de Motor\"]"`
	Cost        *float64            `json:"cost,omitempty" example:"350.50"`
	Attachments []AttachmentRequest `json:"attachments,omitempty"`
}

func (h *MaintenanceHandler) extractUserID(r *http.Request) string {
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
// @Summary      Cadastrar manutenção
// @Description  Registra um novo registro de revisão ou manutenção para o veículo especificado
// @Tags         Maintenances
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID do veículo"
// @Param        payload body CreateMaintenanceRequest true "Dados da manutenção"
// @Success      201 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      403 {object} httpx.ErrorEnvelope
// @Router       /api/v1/cars/{id}/maintenances [post]
func (h *MaintenanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	carID := r.PathValue("id")
	var input CreateMaintenanceRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleMaintenanceError(w, err)
		return
	}
	attachments := make([]domainmaintenance.Attachment, 0, len(input.Attachments))
	for _, att := range input.Attachments {
		attachments = append(attachments, domainmaintenance.Attachment{
			ID:        att.ID,
			Name:      att.Name,
			Size:      att.Size,
			MimeType:  att.MimeType,
			DataUrl:   att.DataUrl,
			CreatedAt: att.CreatedAt,
		})
	}
	maintenance, err := h.service.Create(r.Context(), actorID, carID, maintenusecase.CreateInput{
		Title:       input.Title,
		Description: input.Description,
		Date:        input.Date,
		Mileage:     input.Mileage,
		Types:       input.Types,
		Cost:        input.Cost,
		Attachments: attachments,
	})
	if err != nil {
		handleMaintenanceError(w, err)
		return
	}
	httpx.Created(w, maintenance)
}

// List godoc
// @Summary      Listar manutenções do veículo
// @Description  Retorna o histórico cronológico de manutenções de um veículo
// @Tags         Maintenances
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID do veículo"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      403 {object} httpx.ErrorEnvelope
// @Router       /api/v1/cars/{id}/maintenances [get]
func (h *MaintenanceHandler) List(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	carID := r.PathValue("id")
	maintenances, err := h.service.List(r.Context(), actorID, carID)
	if err != nil {
		handleMaintenanceError(w, err)
		return
	}
	httpx.Success(w, maintenances)
}

// Update godoc
// @Summary      Atualizar manutenção
// @Description  Atualiza dados de uma manutenção existente
// @Tags         Maintenances
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        maintenanceID path string true "ID da manutenção"
// @Param        payload body UpdateMaintenanceRequest true "Dados atualizados"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      403 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/maintenances/{maintenanceID} [put]
func (h *MaintenanceHandler) Update(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	maintID := r.PathValue("maintenanceID")
	var input UpdateMaintenanceRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleMaintenanceError(w, err)
		return
	}
	attachments := make([]domainmaintenance.Attachment, 0, len(input.Attachments))
	for _, att := range input.Attachments {
		attachments = append(attachments, domainmaintenance.Attachment{
			ID:        att.ID,
			Name:      att.Name,
			Size:      att.Size,
			MimeType:  att.MimeType,
			DataUrl:   att.DataUrl,
			CreatedAt: att.CreatedAt,
		})
	}
	maintenance, err := h.service.Update(r.Context(), actorID, maintID, maintenusecase.UpdateInput{
		Title:       input.Title,
		Description: input.Description,
		Date:        input.Date,
		Mileage:     input.Mileage,
		Types:       input.Types,
		Cost:        input.Cost,
		Attachments: attachments,
	})
	if err != nil {
		handleMaintenanceError(w, err)
		return
	}
	httpx.Success(w, maintenance)
}

// Delete godoc
// @Summary      Excluir manutenção
// @Description  Remove o registro de manutenção
// @Tags         Maintenances
// @Produce      json
// @Security     BearerAuth
// @Param        maintenanceID path string true "ID da manutenção"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      403 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/maintenances/{maintenanceID} [delete]
func (h *MaintenanceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	maintID := r.PathValue("maintenanceID")
	if err := h.service.Delete(r.Context(), actorID, maintID); err != nil {
		handleMaintenanceError(w, err)
		return
	}
	httpx.Success(w, map[string]string{"deleted": maintID})
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
