package dto

import "time"

type RegisterRequest struct{
	Email string `json:"email" validate:"required,email"` 
	Password string `json:"password" validate:"required,min=8"` 
}

type LoginRequest struct{
	Email string `json:"email" validate:"required,email"` 
	Password string `json:"password" validate:"required"` 
}

type AuthResponse struct{
	Email string `json:"email"` 
	Token string `json:"token"` 
}

type UserResponse struct {
    ID        int       `json:"id"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}