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

type Repositories struct {
	User
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User: pgdb.NewUserRepository(pg),
	}
}
