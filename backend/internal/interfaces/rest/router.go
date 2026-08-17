package rest

import (
	"net/http"

	portauth "cvmc/internal/application/ports/auth"
	carport "cvmc/internal/application/ports/car"
	maintport "cvmc/internal/application/ports/maintenance"
	userport "cvmc/internal/application/ports/user"
	carusecase "cvmc/internal/application/usecase/car"
	maintusecase "cvmc/internal/application/usecase/maintenance"
	"cvmc/internal/config"
	"cvmc/internal/interfaces/rest/handlers"
	"cvmc/internal/shared/httpx"
	"cvmc/internal/shared/middleware"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(cfg config.Config, users userport.Repository, hasher portauth.PasswordHasher, tokens portauth.TokenService, cars carport.Repository, maintenances maintport.Repository) *Router {
	mux := http.NewServeMux()
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(users, hasher, tokens)
	carHandler := handlers.NewCarHandler(carusecase.NewService(cars, users))
	maintenanceHandler := handlers.NewMaintenanceHandler(maintusecase.NewService(maintenances, cars))

	mux.Handle("GET /health", middleware.Chain(healthHandler.Health(), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("GET /ready", middleware.Chain(healthHandler.Ready(), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("GET /live", middleware.Chain(healthHandler.Live(), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))

	mux.Handle("GET /api/v1", middleware.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.Success(w, map[string]string{"version": "v1"})
	}), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("POST /api/v1/auth/register", middleware.Chain(http.HandlerFunc(authHandler.Register), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("POST /api/v1/auth/login", middleware.Chain(http.HandlerFunc(authHandler.Login), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("POST /api/v1/auth/refresh", middleware.Chain(http.HandlerFunc(authHandler.Refresh), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("GET /api/v1/auth/me", middleware.Chain(http.HandlerFunc(authHandler.Me), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))

	mux.Handle("GET /api/v1/cars", middleware.Chain(http.HandlerFunc(carHandler.List), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("POST /api/v1/cars", middleware.Chain(http.HandlerFunc(carHandler.Create), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("GET /api/v1/cars/{id}", middleware.Chain(http.HandlerFunc(carHandler.Get), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("PUT /api/v1/cars/{id}", middleware.Chain(http.HandlerFunc(carHandler.Update), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("DELETE /api/v1/cars/{id}", middleware.Chain(http.HandlerFunc(carHandler.Delete), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("POST /api/v1/cars/{id}/share", middleware.Chain(http.HandlerFunc(carHandler.Share), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("DELETE /api/v1/cars/{id}/share/{userID}", middleware.Chain(http.HandlerFunc(carHandler.Unshare), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))

	mux.Handle("GET /api/v1/cars/{id}/maintenances", middleware.Chain(http.HandlerFunc(maintenanceHandler.List), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("POST /api/v1/cars/{id}/maintenances", middleware.Chain(http.HandlerFunc(maintenanceHandler.Create), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("PUT /api/v1/maintenances/{maintenanceID}", middleware.Chain(http.HandlerFunc(maintenanceHandler.Update), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("DELETE /api/v1/maintenances/{maintenanceID}", middleware.Chain(http.HandlerFunc(maintenanceHandler.Delete), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))

	return &Router{mux: mux}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
