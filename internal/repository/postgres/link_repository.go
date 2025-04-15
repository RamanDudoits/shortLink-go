package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/RamanDudoits/shortLink-go/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepository struct {
	db *pgxpool.Pool
}

func NewLinkRepository(db *pgxpool.Pool) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) Create(link *repository.Link) (*repository.Link, error) {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var shortLinkID int
	err = tx.QueryRow(ctx,
		`INSERT INTO short_links (link, short_link, clicks, created_at) 
		VALUES ($1, $2, $3, $4) RETURNING id`,
		link.OriginalURL, link.ShortCode, link.ClickCount, time.Now(),
	).Scan(&shortLinkID)
	if err != nil {
		return nil, fmt.Errorf("failed to create short link: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO user_links (user_id, short_link_id) VALUES ($1, $2)`,
		link.UserID, shortLinkID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user link: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	link.ID = shortLinkID
	return link, nil
}

func (r *LinkRepository) FindByID(id int) (*repository.Link, error) {
	query := `
		SELECT sl.id, sl.link, sl.short_link, sl.clicks, sl.created_at, ul.user_id
		FROM short_links sl
		JOIN user_links ul ON sl.id = ul.short_link_id
		WHERE sl.id = $1
	`
	
	var link repository.Link
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&link.ID,
		&link.OriginalURL,
		&link.ShortCode,
		&link.ClickCount,
		&link.CreatedAt,
		&link.UserID,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("link not found")
		}
		return nil, fmt.Errorf("failed to get link: %w", err)
	}
	
	return &link, nil
}

func (r *LinkRepository) FindByUserID(userID int) ([]*repository.Link, error) {
	query := `
		SELECT sl.id, sl.link, sl.short_link, sl.clicks, sl.created_at, ul.user_id
		FROM short_links sl
		JOIN user_links ul ON sl.id = ul.short_link_id
		WHERE ul.user_id = $1
		ORDER BY sl.created_at DESC
	`
	
	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query links: %w", err)
	}
	defer rows.Close()
	
	var links []*repository.Link
	for rows.Next() {
		var link repository.Link
		if err := rows.Scan(
			&link.ID,
			&link.OriginalURL,
			&link.ShortCode,
			&link.ClickCount,
			&link.CreatedAt,
			&link.UserID,
		); err != nil {
			return nil, fmt.Errorf("failed to scan link: %w", err)
		}
		links = append(links, &link)
	}
	
	return links, nil
}

func (r *LinkRepository) FindByURLAndUser(url string, userID int) (*repository.Link, error) {
	query := `
		SELECT sl.id, sl.link, sl.short_link, sl.clicks, sl.created_at, ul.user_id
		FROM short_links sl
		JOIN user_links ul ON sl.id = ul.short_link_id
		WHERE sl.link = $1 AND ul.user_id = $2
		LIMIT 1
	`
	
	var link repository.Link
	err := r.db.QueryRow(context.Background(), query, url, userID).Scan(
		&link.ID,
		&link.OriginalURL,
		&link.ShortCode,
		&link.ClickCount,
		&link.CreatedAt,
		&link.UserID,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find link: %w", err)
	}
	
	return &link, nil
}

func (r *LinkRepository) Update(id int, updates map[string]interface{}) (*repository.Link, error) {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	query := "UPDATE short_links SET "
	params := []interface{}{}
	paramCount := 1

	for field, value := range updates {
		query += fmt.Sprintf("%s = $%d, ", field, paramCount)
		params = append(params, value)
		paramCount++
	}

	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, link, short_link, clicks, created_at", paramCount)
	params = append(params, id)

	var link repository.Link
	err = tx.QueryRow(context.Background(), query, params...).Scan(
		&link.ID,
		&link.OriginalURL,
		&link.ShortCode,
		&link.ClickCount,
		&link.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update link: %w", err)
	}

	err = tx.QueryRow(context.Background(),
		"SELECT user_id FROM user_links WHERE short_link_id = $1", id,
	).Scan(&link.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user_id: %w", err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &link, nil
}

func (r *LinkRepository) Delete(id, userID int) error {
	var exists bool
	err := r.db.QueryRow(context.Background(),
		`SELECT EXISTS(
			SELECT 1 FROM user_links 
			WHERE short_link_id = $1 AND user_id = $2
		)`, id, userID).Scan(&exists)
	
	if err != nil {
		return fmt.Errorf("failed to check link ownership: %w", err)
	}
	
	if !exists {
		return fmt.Errorf("link not found or access denied")
	}

	_, err = r.db.Exec(context.Background(),
		"DELETE FROM short_links WHERE id = $1", id)
	
	if err != nil {
		return fmt.Errorf("failed to delete link: %w", err)
	}
	
	return nil
}

func (r *LinkRepository) Find(filter map[string]interface{}) (*repository.Link, error) {
    query := `
        SELECT sl.id, sl.link, sl.short_link, sl.clicks, sl.created_at, ul.user_id
        FROM short_links sl
        JOIN user_links ul ON sl.id = ul.short_link_id
    `

    var conditions []string
    var args []interface{}
    argNum := 1

    validFields := map[string]string{
        "id":          "sl.id",
        "link":        "sl.link",
        "short_link":  "sl.short_link",
        "user_id":     "ul.user_id",
    }

    for field, value := range filter {
        dbField, ok := validFields[field]
        if !ok {
            continue
        }

        conditions = append(conditions, fmt.Sprintf("%s = $%d", dbField, argNum))
        args = append(args, value)
        argNum++
    }

    if len(conditions) > 0 {
        query += " WHERE " + strings.Join(conditions, " AND ")
    }

    query += " LIMIT 1"

    var link repository.Link
    err := r.db.QueryRow(context.Background(), query, args...).Scan(
        &link.ID,
        &link.OriginalURL,
        &link.ShortCode,
        &link.ClickCount,
        &link.CreatedAt,
        &link.UserID,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to find link: %w", err)
    }

    return &link, nil
}