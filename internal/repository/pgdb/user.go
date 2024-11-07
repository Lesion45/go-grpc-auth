package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"grpc-auth/internal/models"
	"grpc-auth/internal/repository"
	"grpc-auth/pkg/postgres"
)

type UserRepository struct {
	*postgres.Postgres
}

func NewUserRepository(pg *postgres.Postgres) *UserRepository {
	return &UserRepository{pg}
}

func (r *UserRepository) SaveUser(ctx context.Context, email string, passHash []byte) (uuid.UUID, error) {
	const op = "repository.user.SaveUser"

	var userID uuid.UUID

	sault := uniuri.New()

	query := `INSERT INTO users_schema.user(email, password_hash, salt) VALUES(@userEmail, @userPasswordHash, @userSalt) RETURNING id`
	args := pgx.NamedArgs{
		"userEmail":        email,
		"userPasswordHash": passHash,
		"userSalt":         sault,
	}

	err := r.DB.QueryRow(ctx, query, args).Scan(&userID)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok {
			if pgErr.Code == "23505" {
				return uuid.Nil, fmt.Errorf("%s: %w", op, repository.ErrUserExists)
			}
		} else {
			return uuid.Nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return userID, nil
}

func (r *UserRepository) GetUser(ctx context.Context, email string) (models.User, error) {
	const op = "repository.user.GetUser"

	var userID uuid.UUID
	var userEmail string
	var userPasswordHash []byte
	var userSalt string

	query := `SELECT id, email, password_hash FROM users_schema.user WHERE email = @userEmail`
	args := pgx.NamedArgs{
		"userEmail": email,
	}

	err := r.DB.QueryRow(ctx, query, args).Scan(&userID, &userEmail, &userPasswordHash, &userSalt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		} else {
			return models.User{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	user := models.User{
		ID:           userID,
		Email:        userEmail,
		PasswordHash: userPasswordHash,
		Salt:         userSalt,
	}

	return user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, email string) (uuid.UUID, error) {
	const op = "repository.user.DeleteUser"

	var deletedUserID uuid.UUID

	query := `DELETE FROM users_schema.user WHERE email = @userEmail RETURNING id`
	args := pgx.NamedArgs{
		"userEmail": email,
	}

	err := r.DB.QueryRow(ctx, query, args).Scan(&deletedUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		} else {
			return uuid.Nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return deletedUserID, nil
}
