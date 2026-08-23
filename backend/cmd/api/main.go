package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cvmc/internal/bootstrap"
	"cvmc/internal/config"
)

// @title           CVMC API
// @version         1.0
// @description     API REST para gestão inteligente de veículos, revisões e manutenções do CVMC (Como Vai Meu Carro).
// @termsOfService  http://swagger.io/terms/

// @contact.name   Suporte CVMC
// @contact.email  suporte@cvmc.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Digite "Bearer " seguido do token JWT obtido no login.
func main() {
	cfg := config.Load()

	initCtx, initCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer initCancel()

	app, err := bootstrap.New(initCtx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	srv := &http.Server{
		Addr:              cfg.HTTPAddress(),
		Handler:           app.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	go func() {
		log.Printf("server starting on %s", cfg.HTTPAddress())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	if err := app.Close(ctx); err != nil {
		log.Printf("error closing database connection: %v", err)
	}
}
