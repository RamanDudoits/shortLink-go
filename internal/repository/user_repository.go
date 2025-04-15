package repository

import (
	"context"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetAllUsers(ctx context.Context) ([]*User, error)
}

type User struct {
	ID           int
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}