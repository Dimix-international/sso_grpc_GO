package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"syscall"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main() {
	//инициализируем конфиг
	cfg := config.MustLoad()
	
	//инициализируем логгер
	log := setupLogger(cfg.Env)

	log.Info("starting app", slog.Any("env", cfg))

	//инициализация и запуск gRPC сервера
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	go application.GRPCSrv.MustRun()

	//Graceful shutdown - остановка процессов

	stop := make(chan os.Signal, 1)
	//слушаем сигналы операционной системы, получим сигнал о завершении  и запишем в канал
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	//чтение из канала блокирующая операция, пока не придет запись, код дальше не идет
	sign := <- stop

	//у бд как Postgres так же должны быть методы для Graceful shutdown, для workers и т.д.
	
	application.GRPCSrv.Stop()
	log.Info("application stopped with signal: ", sign)

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal: 
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}), //LevelDebug - все все логи
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}