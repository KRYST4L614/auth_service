package grpc

import (
	"fmt"
	authgrpc "github.com/KRYST4L614/auth_service/internal/grpc/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func NewApp(
	log *slog.Logger,
	authService authgrpc.Auth,
	port string,
) *App {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer, authService)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port))

	lis, err := net.Listen("tcp", ":"+a.port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server started", slog.String("addr", lis.Addr().String()))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(
		slog.String("op", op)).Info("stopping grpc server")

	a.gRPCServer.GracefulStop()
}
