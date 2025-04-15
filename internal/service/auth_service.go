package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RamanDudoits/shortLink-go/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
    Register(email, password string) (*User, error)
    Authenticate(email, password string) (*User, error)
	GenerateToken(userID int, email string) (string, error)
	ValidateToken(tokenString string) (int, string, error)
	GetAllUsers() ([]*User, error)
}

type User struct {
	ID       int
	Email    string
	Password string
	CreatedAt time.Time
}

type AuthService struct {
	repo repository.UserRepository
	jwtSecret string
	tokenExpiry time.Duration
}

func NewAuthService(repo repository.UserRepository, jwtSecret string, tokenExpiry time.Duration) *AuthService {
	return &AuthService{
		repo: repo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

func (s *AuthService) Register(email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.repo.CreateUser(context.Background(), email, string(hashedPassword))

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	return &User{
		ID: user.ID,
		Email: user.Email,
	}, nil
}

func (s *AuthService) Authenticate(email, password string) (*User, error) {
	user, err := s.repo.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return &User{
		ID:       user.ID,
		Email:    user.Email,
	}, nil
}

func (s *AuthService) GenerateToken(userID int, email string) (string, error){
	claims := jwt.MapClaims{
		"sub": userID,
		"email": email,
		"exp": time.Now().Add(s.tokenExpiry).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (int, string, error){
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, "", fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["sub"].(float64))
		email := claims["email"].(string)
		return userID, email, nil
	}

	return 0, "", errors.New("invalid token")
}

func (s *AuthService) GetAllUsers() ([]*User, error) {
    dbUsers, err := s.repo.GetAllUsers(context.Background())
    if err != nil {
        return nil, fmt.Errorf("failed to get users 2: %w", err)
    }

    users := make([]*User, len(dbUsers))
    for i, dbUser := range dbUsers {
        users[i] = &User{
            ID:        dbUser.ID,
            Email:     dbUser.Email,
            CreatedAt: dbUser.CreatedAt,
        }
    }

    return users, nil
}
