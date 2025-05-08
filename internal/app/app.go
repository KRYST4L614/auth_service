package app

import (
	grpcapp "github.com/KRYST4L614/auth_service/internal/app/grpc"
	"github.com/KRYST4L614/auth_service/internal/services/auth"
	"github.com/KRYST4L614/auth_service/internal/storage/sqlite"
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
	storage, err := sqlite.NewStorage(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.NewApp(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
