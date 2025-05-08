package suite

import (
	"context"
	"github.com/KRYST4L614/auth_service/internal/config"
	ssov1 "github.com/KRYST4L614/auth_service_protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
)

const grpcHost = "localhost"

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	path := "../config/local.yaml"

	cfg := config.MustLoadByPath(path)

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	clientConn, err := grpc.NewClient(
		net.JoinHostPort(grpcHost, cfg.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(clientConn),
	}
}
