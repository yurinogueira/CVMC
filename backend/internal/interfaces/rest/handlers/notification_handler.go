package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	portauth "cvmc/internal/application/ports/auth"
	notificationusecase "cvmc/internal/application/usecase/notification"
	"cvmc/internal/shared/httpx"
)

type NotificationHandler struct {
	service *notificationusecase.Service
	tokens  portauth.TokenService
}

func NewNotificationHandler(service *notificationusecase.Service, tokens ...portauth.TokenService) *NotificationHandler {
	var tokenService portauth.TokenService
	if len(tokens) > 0 {
		tokenService = tokens[0]
	}
	return &NotificationHandler{
		service: service,
		tokens:  tokenService,
	}
}

type RegisterDeviceTokenRequest struct {
	Token    string `json:"token" example:"fcm_device_token_string_123"`
	Platform string `json:"platform" example:"android"`
}

type UnregisterDeviceTokenRequest struct {
	Token string `json:"token" example:"fcm_device_token_string_123"`
}

func (h *NotificationHandler) extractUserID(r *http.Request) string {
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

// RegisterDeviceToken godoc
// @Summary      Registrar token FCM do dispositivo
// @Description  Vincula ou atualiza o token push FCM do aparelho para o usuário autenticado
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body RegisterDeviceTokenRequest true "Token e plataforma do dispositivo"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/devices/token [post]
func (h *NotificationHandler) RegisterDeviceToken(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var input RegisterDeviceTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid payload", nil)
		return
	}

	dt, err := h.service.RegisterDeviceToken(r.Context(), actorID, input.Token, input.Platform)
	if err != nil {
		handleNotificationError(w, err)
		return
	}
	httpx.Success(w, dt)
}

// UnregisterDeviceToken godoc
// @Summary      Remover token FCM do dispositivo
// @Description  Desvincula o token push FCM do aparelho ao efetuar logout
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body UnregisterDeviceTokenRequest true "Token do dispositivo a remover"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Router       /api/v1/devices/token [delete]
func (h *NotificationHandler) UnregisterDeviceToken(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var input UnregisterDeviceTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid payload", nil)
		return
	}

	if err := h.service.UnregisterDeviceToken(r.Context(), actorID, input.Token); err != nil {
		handleNotificationError(w, err)
		return
	}
	httpx.Success(w, map[string]string{"message": "Token do dispositivo removido com sucesso"})
}

// SendTestNotification godoc
// @Summary      Disparar notificação push de teste
// @Description  Envia uma notificação push de teste e verifica alertas preventivos para os dispositivos registrados do usuário autenticado
// @Tags         Notifications
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/notifications/test [post]
func (h *NotificationHandler) SendTestNotification(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	sentCount, err := h.service.SendTestNotification(r.Context(), actorID)
	if err != nil {
		handleNotificationError(w, err)
		return
	}

	maintenanceAlerts, _ := h.service.CheckMaintenanceAlerts(r.Context(), actorID)
	httpx.Success(w, map[string]any{
		"message":           "Notificação push de teste enviada com sucesso!",
		"devicesNotified":   sentCount,
		"maintenanceAlerts": maintenanceAlerts,
	})
}

func handleNotificationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, notificationusecase.ErrInvalidPayload):
		httpx.Error(w, http.StatusBadRequest, "Dados do dispositivo inválidos", nil)
	case errors.Is(err, notificationusecase.ErrNoDevicesRegistered):
		httpx.Error(w, http.StatusNotFound, "Nenhum dispositivo móvel registrado para receber notificações push", nil)
	case errors.Is(err, notificationusecase.ErrDeliveryFailed):
		httpx.Error(w, http.StatusBadGateway, "Falha ao entregar notificação push via FCM", nil)
	default:
		httpx.Error(w, http.StatusInternalServerError, "Unexpected error", nil)
	}
}
