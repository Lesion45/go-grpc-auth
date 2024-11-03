package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"grpc-auth/config"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func (pg *Postgres) PostgresHealthCheck(ctx context.Context) error {
	if err := pg.DB.Ping(ctx); err != nil {
		return err
	}
	return nil
}

var (
	pgInstance *Postgres
)

func NewPG(ctx context.Context, storage config.Storage) (*Postgres, error) {
	const op = "storage.postgres.NewPG"

	DSN := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		storage.Host, storage.Port, storage.User, storage.Password, storage.Name)

	db, err := pgxpool.New(ctx, DSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pgInstance = &Postgres{DB: db}

	return pgInstance, nil
}
