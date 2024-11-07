package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"grpc-auth/internal/lib/jwt"
	"grpc-auth/internal/repository"
	"grpc-auth/pkg/logger/sl"
	"log/slog"
	"time"
)

type Auth interface {
	Login(ctx context.Context, email string, password string) (string, error)
	RegisterNewUser(ctx context.Context, email string, password string) (string, error)
}

type AuthDependencies struct {
	Log      *slog.Logger
	Repos    *repository.Repositories
	TokenTTL time.Duration
}

type AuthService struct {
	log            *slog.Logger
	userRepository repository.User
	tokenTTL       time.Duration
}

func New(deps AuthDependencies) *AuthService {
	return &AuthService{
		log:            deps.Log,
		userRepository: deps.Repos.User,
		tokenTTL:       deps.TokenTTL,
	}
}

func (a *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "service.Auth.Login"

	a.log.With(slog.String("op", op))
	a.log.Info("attempting to login user")

	user, err := a.userRepository.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, err)
		}

		a.log.Error("failed to get user", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	a.log.Info("user logged successfully")

	token, err := jwt.NewToken(user, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *AuthService) RegisterNewUser(ctx context.Context, email string, password string) (string, error) {
	const op = "service.Auth.Registration"

	a.log.With(slog.String("op", op))
	a.log.Info("attempting to register new user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("failed to generate password hash", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userRepository.SaveUser(ctx, email, passwordHash)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			a.log.Warn("user already exists", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, err)
		}
		a.log.Error("failed to save user", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	userID := id.String()
	if userID == "" {
		a.log.Error("failed to save user", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("user successfully saved")

	return userID, nil
}
