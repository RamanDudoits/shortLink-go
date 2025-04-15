package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/RamanDudoits/shortLink-go/internal/transport/http/handler"
)

func NewRouter(authHandler handler.AuthHandlerInterface, linkHandler handler.LinkHandlerInterface) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	}))

	r.Group(func(r chi.Router) {
		r.Post("/api/auth/login", authHandler.Login)
		r.Post("/api/auth/register", authHandler.Register)
		r.Get("/{shortLink}", linkHandler.Redirect)
		r.Get("/api/admin/users", authHandler.GetAllUsers)
	})

	r.Group(func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)

		r.Get("/api/links", linkHandler.List)
		r.Post("/api/links", linkHandler.Store)
		r.Delete("/api/links/destroy", linkHandler.Destroy)
		r.Get("/api/links/{shortLink}", linkHandler.Get)
		r.Patch("/api/links/{shortLink}/update", linkHandler.Update)
	})
	return r
}