package bootstrap

import (
	"cvmc/internal/config"
	bcryptinfra "cvmc/internal/infrastructure/auth/bcrypt"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	carrepo "cvmc/internal/infrastructure/car/memory"
	maintrepo "cvmc/internal/infrastructure/maintenance/memory"
	memoryuser "cvmc/internal/infrastructure/user/memory"
	"cvmc/internal/interfaces/rest"
)

type App struct {
	handler *rest.Router
}

func New(cfg config.Config) *App {
	users := memoryuser.NewRepository()
	hasher := bcryptinfra.NewHasher()
	tokens := jwtauth.NewProvider(cfg.JWTSecret, cfg.JWTRefreshSecret)
	cars := carrepo.NewRepository()
	maintenances := maintrepo.NewRepository()
	return &App{handler: rest.NewRouter(cfg, users, hasher, tokens, cars, maintenances)}
}

func (a *App) Handler() *rest.Router {
	return a.handler
}
