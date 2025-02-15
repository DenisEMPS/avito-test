package main

import (
	"avito/internal/config"
	"avito/internal/handler"
	"avito/internal/repository"
	"avito/internal/service"
	"avito/server"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.Any("cfg", cfg))

	db, err := repository.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Error("failed connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	repos := repository.NewRepositry(db, log)
	service := service.NewService(repos, log)
	handler := handler.NewHandler(service, log)

	server := new(server.Server)

	go func() {
		if err := server.Run(cfg.Server, handler.InitRoutes()); err != nil {
			log.Error("failed to run server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	log.Info("App started", slog.Any("port", cfg.Server.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := server.Shutdown(context.Background()); err != nil {
		log.Error("error ocured on server shutting down", slog.Any("error", err))
	}
	if err := db.Close(); err != nil {
		log.Info("error ocured on db connection close", slog.Any("error", err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
