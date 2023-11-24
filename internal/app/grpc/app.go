package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	authgrpc "sso/internal/grpc/auth"

	"google.golang.org/grpc"
)

type App struct {
	log *slog.Logger
	gRPCServer *grpc.Server
	port int
}

//конструктор
func New(log *slog.Logger, authService authgrpc.Auth, port int) *App {
	//запуск и конфигурация grpc сервера
	gRPCServer := grpc.NewServer()

	//подключаем обработчик
	authgrpc.Register(gRPCServer, authService)

	return &App {
		log: log,
		gRPCServer: gRPCServer,
		port: port,
	}
}
//Must функции - вместо возврата ошибки - panic
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

//функция запуска grpc сервера
func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port)) //чтобы в логах была инфа об op

	//слушатель для прослушки tcp сообщений
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	//адрресс по которому обрабатывается tcp соединение
	log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	//запуск сервера
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

//остановка gRPC сервера
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op), slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}