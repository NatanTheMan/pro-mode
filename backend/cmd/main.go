package main

import (
	"context"
	"fmt"
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
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	repo := repository.NewUserRepository(pool)
	svc := service.NewUserService(repo, cfg.AdminID, cfg.DevID)

	err = telegram.Run(ctx, cfg, svc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
