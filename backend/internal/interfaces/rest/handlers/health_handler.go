package handlers

import (
	"net/http"

	"cvmc/internal/shared/httpx"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.Success(w, map[string]string{"status": "ok"})
	})
}

func (h *HealthHandler) Ready() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.Success(w, map[string]string{"status": "ready"})
	})
}

func (h *HealthHandler) Live() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.Success(w, map[string]string{"status": "alive"})
	})
}
