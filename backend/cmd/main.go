package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/NatanTheMan/pro-mode/internal/config"
	"github.com/NatanTheMan/pro-mode/internal/database"
	"github.com/NatanTheMan/pro-mode/internal/repository"
	"github.com/NatanTheMan/pro-mode/internal/service"
	"github.com/NatanTheMan/pro-mode/internal/telegram"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		slog.Error("bot error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := newLogger()
	slog.SetDefault(logger)

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}

	userRepo := repository.NewUserRepository(pool)
	userSvc := service.NewUserService(userRepo, cfg.AdminID, cfg.DevID)
	categoryRepo := repository.NewCategoryRepository(pool)
	categorySvc := service.NewCategoryService(categoryRepo, cfg.AdminID, cfg.DevID)

	err = telegram.Run(ctx, cfg, logger, userSvc, categorySvc)
	if err != nil {
		return err
	}
	logger.Info("bot starting")

	return nil
}

func newLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}
