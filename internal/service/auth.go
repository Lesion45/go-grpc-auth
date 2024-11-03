package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"grpc-auth/internal/lib/jwt"
	"grpc-auth/internal/repository"
	"grpc-auth/pkg/logger/sl"
	"log/slog"
	"time"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID string) (string, error)
	RegisterNewUser(ctx context.Context, email string, password string) (string, error)
	RegisterNewApp(ctx context.Context, name string, appSecret string, adminKey string) (string, error)
}

type AuthDependencies struct {
	Log      *slog.Logger
	Repos    *repository.Repositories
	TokenTTL time.Duration
}

type AuthService struct {
	log            *slog.Logger
	userRepository repository.User
	appRepository  repository.App
	tokenTTL       time.Duration
}

func New(deps AuthDependencies) *AuthService {
	return &AuthService{
		log:            deps.Log,
		userRepository: deps.Repos.User,
		appRepository:  deps.Repos.App,
		tokenTTL:       deps.TokenTTL,
	}
}

func (a *AuthService) Login(ctx context.Context, email string, password string, appID string) (string, error) {
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

	id, err := uuid.Parse(appID)
	if err != nil {
		a.log.Error("failed to get app", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.appRepository.GetApp(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAppNotFound) {
			a.log.Warn("app not found", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidData)
		}

		a.log.Error("failed to get app", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
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

func (a *AuthService) RegisterNewApp(ctx context.Context, name string, secretKey string, adminKey string) (string, error) {
	const op = "service.Auth.AddApp"

	a.log.With(slog.String("op", op))
	a.log.Info("attempting to add app")

	id, err := a.appRepository.SaveApp(ctx, name, secretKey)
	if err != nil {
		if errors.Is(err, repository.ErrAppExists) {
			a.log.Warn("app already exists")

			return "", fmt.Errorf("%s: %w", op, ErrInvalidData)
		} else {
			a.log.Error("failed to add app", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	appID := id.String()
	if appID == "" {
		a.log.Error("failed to add app", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return appID, nil
}

func (a *AuthService) DeleteApp(ctx context.Context, appID string) error {
	const op = "service.Auth.DeleteApp"

	a.log.With(slog.String("op", op))
	a.log.Info("attempting to delete app")

	id, err := uuid.Parse(appID)
	if err != nil {
		a.log.Error("failed to delete app", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	err = a.appRepository.DeleteApp(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAppNotFound) {
			a.log.Warn("app not found")

			return fmt.Errorf("%s: %w", op, err)
		} else {
			a.log.Error("failed to delete app", sl.Err(err))

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
