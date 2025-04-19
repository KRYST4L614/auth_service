package app

import (
	grpcapp "github.com/KRYST4L614/auth_service/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(
	log *slog.Logger,
	grpcPort string,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO: инициализировать хранилище

	// TODO: init auth service

	grpcApp := grpcapp.NewApp(log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
