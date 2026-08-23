package handlers

import (
	"net/http"

	"cvmc/internal/shared/httpx"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health godoc
// @Summary      Health check
// @Description  Retorna o status geral de saúde da aplicação
// @Tags         Health
// @Produce      json
// @Success      200 {object} httpx.SuccessEnvelope
// @Router       /health [get]
func (h *HealthHandler) Health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.Success(w, map[string]string{"status": "ok"})
	})
}

// Ready godoc
// @Summary      Readiness probe
// @Description  Verifica se a aplicação está pronta para receber tráfego
// @Tags         Health
// @Produce      json
// @Success      200 {object} httpx.SuccessEnvelope
// @Router       /ready [get]
func (h *HealthHandler) Ready() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.Success(w, map[string]string{"status": "ready"})
	})
}

// Live godoc
// @Summary      Liveness probe
// @Description  Verifica se a aplicação está viva em execução
// @Tags         Health
// @Produce      json
// @Success      200 {object} httpx.SuccessEnvelope
// @Router       /live [get]
func (h *HealthHandler) Live() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.Success(w, map[string]string{"status": "alive"})
	})
}
