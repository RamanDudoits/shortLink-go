package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)


type Postgres struct {
	Poll *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("enaible to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("enaible to ping database: %w", err)	
	}

	return &Postgres{
		Poll: pool,
	}, nil
}

func (p *Postgres) Close() {
	p.Poll.Close()
}