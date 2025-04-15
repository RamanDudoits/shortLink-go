package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/RamanDudoits/shortLink-go/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, email, passwordHash string) (*repository.User, error) {
	var user repository.User
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash, created_at
	`

	err := r.db.QueryRow(ctx, query, email, passwordHash).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, errors.New("user already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	var user repository.User
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*repository.User, error) {
	var user repository.User
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func isDuplicateKeyError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*repository.User, error) {
    query := `
        SELECT id, email, created_at
        FROM users
        ORDER BY created_at DESC
    `

    rows, err := r.db.Query(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get users 3: %w", err)
    }
    defer rows.Close()

    var users []*repository.User
    for rows.Next() {
        var user repository.User
        if err := rows.Scan(
            &user.ID,
            &user.Email,
            &user.CreatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan user: %w", err)
        }
        users = append(users, &user)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }

    return users, nil
}