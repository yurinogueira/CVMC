package rest

import (
	"net/http"

	_ "cvmc/docs"
	portauth "cvmc/internal/application/ports/auth"
	carport "cvmc/internal/application/ports/car"
	emailport "cvmc/internal/application/ports/email"
	maintport "cvmc/internal/application/ports/maintenance"
	userport "cvmc/internal/application/ports/user"
	carusecase "cvmc/internal/application/usecase/car"
	fipeusecase "cvmc/internal/application/usecase/fipe"
	maintusecase "cvmc/internal/application/usecase/maintenance"
	"cvmc/internal/config"
	"cvmc/internal/interfaces/rest/handlers"
	"cvmc/internal/shared/httpx"
	"cvmc/internal/shared/middleware"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Router struct {
	handler http.Handler
}

func NewRouter(cfg config.Config, users userport.Repository, hasher portauth.PasswordHasher, tokens portauth.TokenService, emailSender emailport.Sender, cars carport.Repository, maintenances maintport.Repository, fipeService *fipeusecase.Service) *Router {
	mux := http.NewServeMux()
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(users, hasher, tokens, emailSender, cfg.CookieDomain, cfg.CookieSecure)
	userHandler := handlers.NewUserHandler(users, cars, hasher, tokens)
	carHandler := handlers.NewCarHandler(carusecase.NewService(cars, users), tokens)
	maintenanceHandler := handlers.NewMaintenanceHandler(maintusecase.NewService(maintenances, cars), tokens)
	fipeHandler := handlers.NewFipeHandler(fipeService, tokens)

	// Rate limiters
	globalLimiter := middleware.NewRateLimiter(100.0/60.0, 100, cfg.TrustedProxies...) // 100 req/min burst 100
	authLimiter := middleware.NewRateLimiter(10.0/60.0, 10, cfg.TrustedProxies...)     // 10 req/min burst 10
	strictLimiter := middleware.NewRateLimiter(5.0/60.0, 5, cfg.TrustedProxies...)     // 5 req/min burst 5

	// Swagger: only available in debug mode
	if cfg.LogLevel == "debug" {
		mux.Handle("/swagger/", httpSwagger.WrapHandler)
	}

	mux.Handle("GET /health", middleware.Chain(healthHandler.Health(), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("GET /ready", middleware.Chain(healthHandler.Ready(), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("GET /live", middleware.Chain(healthHandler.Live(), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))

	mux.Handle("GET /api/v1", middleware.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.Success(w, map[string]string{"version": "v1"})
	}), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))

	// Auth endpoints with stricter rate limiting
	authChain := func(h http.HandlerFunc) http.Handler {
		return authLimiter.Limit(middleware.Chain(h, middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	}
	strictChain := func(h http.HandlerFunc) http.Handler {
		return strictLimiter.Limit(middleware.Chain(h, middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	}

	mux.Handle("POST /api/v1/auth/register", authChain(authHandler.Register))
	mux.Handle("POST /api/v1/auth/login", authChain(authHandler.Login))
	mux.Handle("POST /api/v1/auth/refresh", middleware.Chain(http.HandlerFunc(authHandler.Refresh), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("GET /api/v1/auth/me", middleware.Chain(http.HandlerFunc(authHandler.Me), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("POST /api/v1/auth/logout", middleware.Chain(http.HandlerFunc(authHandler.Logout), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("POST /api/v1/auth/forgot-password", strictChain(authHandler.ForgotPassword))
	mux.Handle("POST /api/v1/auth/reset-password", strictChain(authHandler.ResetPassword))
	mux.Handle("POST /api/v1/auth/verify-email", authChain(authHandler.VerifyEmail))
	mux.Handle("POST /api/v1/auth/resend-verification", strictChain(authHandler.ResendVerification))

	// User & Profile endpoints
	mux.Handle("GET /api/v1/user/profile", middleware.Chain(http.HandlerFunc(userHandler.GetProfile), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("PUT /api/v1/user/profile", middleware.Chain(http.HandlerFunc(userHandler.UpdateProfile), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("PUT /api/v1/user/password", middleware.Chain(http.HandlerFunc(userHandler.UpdatePassword), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))

	// Cars endpoints
	mux.Handle("GET /api/v1/cars", middleware.Chain(http.HandlerFunc(carHandler.List), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("POST /api/v1/cars", middleware.Chain(http.HandlerFunc(carHandler.Create), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("GET /api/v1/cars/{id}", middleware.Chain(http.HandlerFunc(carHandler.Get), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("PUT /api/v1/cars/{id}", middleware.Chain(http.HandlerFunc(carHandler.Update), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("DELETE /api/v1/cars/{id}", middleware.Chain(http.HandlerFunc(carHandler.Delete), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("POST /api/v1/cars/{id}/share", middleware.Chain(http.HandlerFunc(carHandler.Share), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("DELETE /api/v1/cars/{id}/share/{userID}", middleware.Chain(http.HandlerFunc(carHandler.Unshare), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))

	// Maintenance endpoints
	mux.Handle("GET /api/v1/cars/{id}/maintenances", middleware.Chain(http.HandlerFunc(maintenanceHandler.List), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("POST /api/v1/cars/{id}/maintenances", middleware.Chain(http.HandlerFunc(maintenanceHandler.Create), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("PUT /api/v1/maintenances/{maintenanceID}", middleware.Chain(http.HandlerFunc(maintenanceHandler.Update), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("DELETE /api/v1/maintenances/{maintenanceID}", middleware.Chain(http.HandlerFunc(maintenanceHandler.Delete), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))

	// Fipe endpoints
	mux.Handle("GET /api/v1/fipe/{vehicleType}/brands", middleware.Chain(http.HandlerFunc(fipeHandler.ListBrands), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("GET /api/v1/fipe/{vehicleType}/brands/{brandId}/models", middleware.Chain(http.HandlerFunc(fipeHandler.ListModels), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("GET /api/v1/fipe/{vehicleType}/brands/{brandId}/models/{modelId}/years", middleware.Chain(http.HandlerFunc(fipeHandler.ListYears), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))
	mux.Handle("GET /api/v1/fipe/{vehicleType}/brands/{brandId}/models/{modelId}/years/{yearId}", middleware.Chain(http.HandlerFunc(fipeHandler.GetVehicleDetail), middleware.RequestID, middleware.StructuredLogging(cfg.LogLevel)))

	// Apply global middleware: Security Headers → CORS → Rate Limit → Body Limit → Mux
	handler := middleware.SecurityHeaders(
		middleware.CORS(
			globalLimiter.Limit(
				middleware.BodyLimit(1<<20)(mux), // 1MB body limit
			),
			cfg.AllowedOrigins,
		),
	)

	return &Router{handler: handler}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}
