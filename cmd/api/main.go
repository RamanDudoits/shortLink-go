package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/RamanDudoits/shortLink-go/internal/config"
	"github.com/RamanDudoits/shortLink-go/internal/repository/postgres"
	"github.com/RamanDudoits/shortLink-go/internal/service"
	router "github.com/RamanDudoits/shortLink-go/internal/transport/http"
	"github.com/RamanDudoits/shortLink-go/internal/transport/http/handler"
	"github.com/RamanDudoits/shortLink-go/pkg/validator"
)


func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := postgres.New(ctx, cfg.DB.DSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db.Poll)
	linkRepo := postgres.NewLinkRepository(db.Poll)

	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.Expire,)
	linkService := service.NewLinkService(linkRepo)

	validator := validator.New()

	authHandler := handler.NewAuthHandler(authService, validator)
	linkHandler := handler.NewLinkHandler(linkService, linkService)

	r := router.NewRouter(authHandler, linkHandler)
	
	server := &http.Server{
		Addr: cfg.HTTP.Addr,
		Handler: r,
		ReadTimeout: cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	go func() {
		log.Printf("Server started on %s", cfg.HTTP.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()


	<-ctx.Done()
	log.Println("shutting down server...")
}