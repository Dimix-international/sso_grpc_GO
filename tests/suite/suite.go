package suite

import (
	"context"
	"net"
	"os"
	"sso/internal/config"
	"strconv"
	"testing"

	ssov1 "github.com/Dimix-international/auth_protoc/gen/go/sso"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T                  // Потребуется для вызова методов *testing.T внутри Suite
	Cfg        *config.Config   // Конфигурация приложения
	AuthClient ssov1.AuthClient // Клиент для взаимодействия с gRPC-сервером
}

const (
	grpcHost = "localhost"
)

// New creates new test suite.
//
// TODO: for pipeline tests we need to wait for app is ready
func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath(configPath())

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	//создаем клиента
	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials())) // Используем insecure-коннект для тестов
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func configPath() string {
	const key = "CONFIG_LOCAL_TEST_PATH"

	godotenv.Load()

	if v := os.Getenv(key); v != "" {
		return v
	}

	return "../config/local_tests.yaml"
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port)) //JoinHostPort - объединяет хост и порт в один адресс
}