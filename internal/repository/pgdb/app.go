package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"grpc-auth/internal/models"
	"grpc-auth/internal/repository"
	"grpc-auth/pkg/postgres"
)

type AppRepository struct {
	*postgres.Postgres
}

func NewAppRepository(pg *postgres.Postgres) *AppRepository {
	return &AppRepository{pg}
}

func (r *AppRepository) SaveApp(ctx context.Context, name string, secret string) (uuid.UUID, error) {
	const op = "repository.app.SaveApp"

	var appID uuid.UUID

	query := `INSERT INTO apps_schema.app(name, secret) VALUES (@appName, @appSecret) RETURNING id`
	args := pgx.NamedArgs{
		"appName":   name,
		"appSecret": secret,
	}

	err := r.DB.QueryRow(ctx, query, args).Scan(&appID)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok {
			if pgErr.Code == "23505" {
				return uuid.Nil, fmt.Errorf("%s: %w", op, repository.ErrAppExists)
			}
		} else {
			return uuid.Nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return appID, nil
}

func (r *AppRepository) GetApp(ctx context.Context, id uuid.UUID) (models.App, error) {
	const op = "repository.app.GetApp"

	var appID uuid.UUID
	var appName string
	var secret string

	query := `SELECT id, name, secret FROM apps_schema.app WHERE id = @appID`
	args := pgx.NamedArgs{
		"appID": id,
	}

	err := r.DB.QueryRow(ctx, query, args).Scan(&appID, &appName, &secret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, repository.ErrAppNotFound)
		} else {
			return models.App{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	app := models.App{
		appID,
		appName,
		secret,
	}

	return app, nil
}

func (r *AppRepository) DeleteApp(ctx context.Context, id uuid.UUID) error {
	const op = "repository.app.DeleteApp"

	query := `DELETE FROM apps_schema.app WHERE id = @appID`
	args := pgx.NamedArgs{
		"appID": id,
	}

	commandTag, err := r.DB.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, repository.ErrAppNotFound)
	}

	return nil
}
