package app

import (
	"context"
	"grpc-auth/config"
	grpcapp "grpc-auth/internal/app/grpc"
	"grpc-auth/internal/repository"
	"grpc-auth/internal/service"
	"grpc-auth/pkg/postgres"
	"log/slog"
	"time"
)

type App struct {
	GRPCApp *grpcapp.GRPCApp
}

func New(log *slog.Logger, storageConfig config.Storage, tokenTTL time.Duration, port int) *App {
	pg, err := postgres.NewPG(context.Background(), storageConfig)
	if err != nil {
		panic(err)
	}

	repositories := repository.NewRepositories(pg)

	authService := service.New(service.AuthDependencies{
		Log:      log,
		Repos:    repositories,
		TokenTTL: tokenTTL,
	})

	grpcApp := grpcapp.New(log, authService, port)

	return &App{
		GRPCApp: grpcApp,
	}
}
