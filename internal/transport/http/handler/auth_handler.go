package handler

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/RamanDudoits/shortLink-go/internal/service"
	"github.com/RamanDudoits/shortLink-go/internal/transport/http/dto"
	"github.com/RamanDudoits/shortLink-go/pkg/validator"
	"github.com/go-chi/render"
)

type AuthHandlerInterface interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	AuthMiddleware(next http.Handler) http.Handler
}

type AuthHandler struct {
	authService service.AuthServiceInterface
	validator *validator.Validator
}

func NewAuthHandler(authService service.AuthServiceInterface, validator *validator.Validator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator: validator,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

	if err := h.validator.Validate(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(req.Email, req.Password)

	if err != nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, map[string]string{"error": "User already exists"})
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.Email)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to generate token"})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, dto.AuthResponse{
		Token: token,
		Email: user.Email,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

	if err := h.validator.Validate(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	user, err := h.authService.Authenticate(req.Email, req.Password)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "Invalid credentials"})
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to generate token"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, dto.AuthResponse{
		Token: token,
		Email: user.Email,
	})
}

func (h *AuthHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.authService.GetAllUsers()
    if err != nil {
        log.Printf("GetAllUsers error: %v", err)
        
        render.Status(r, http.StatusInternalServerError)
        render.JSON(w, r, map[string]string{
            "error": "failed to get users",
            "details": err.Error(),
        })
        return
    }

    response := make([]dto.UserResponse, len(users))
    for i, user := range users {
        response[i] = dto.UserResponse{
            ID:        user.ID,
            Email:     user.Email,
            CreatedAt: user.CreatedAt,
        }
    }

    render.Status(r, http.StatusOK)
    render.JSON(w, r, response)
}

func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := extractToken(r)
        if tokenString == "" {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte(`{"error": "Authorization header required"}`))
            return
        }

        userID, _, err := h.authService.ValidateToken(tokenString)
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte(`{"error": "Invalid token"}`))
            return
        }

        ctx := context.WithValue(r.Context(), "userID", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func extractToken(r *http.Request) string {
    bearer := r.Header.Get("Authorization")
    if len(bearer) > 7 && strings.HasPrefix(bearer, "Bearer ") {
        return bearer[7:]
    }
    return ""
}