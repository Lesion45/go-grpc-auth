package repository

import (
	"context"
	"github.com/google/uuid"
	"grpc-auth/internal/models"
	"grpc-auth/internal/repository/pgdb"
	"grpc-auth/pkg/postgres"
)

type User interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uuid.UUID, error)
	GetUser(ctx context.Context, email string) (models.User, error)
	DeleteUser(ctx context.Context, email string) (uuid.UUID, error)
}

type App interface {
	SaveApp(ctx context.Context, name string, secret string) (uuid.UUID, error)
	GetApp(ctx context.Context, id uuid.UUID) (models.App, error)
	DeleteApp(ctx context.Context, id uuid.UUID) error
}

type Repositories struct {
	User
	App
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User: pgdb.NewUserRepository(pg),
		App:  pgdb.NewAppRepository(pg),
	}
}
