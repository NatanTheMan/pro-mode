package telegram

import (
	"context"
	"log/slog"

	"github.com/NatanTheMan/pro-mode/internal/config"
	"github.com/NatanTheMan/pro-mode/internal/service"
	"github.com/go-telegram/bot"
)

func Run(
	ctx context.Context,
	cfg *config.Config,
	logger *slog.Logger,
	userSvc *service.UserService,
	categorySvc *service.CategoryService,
) error {
	opts := []bot.Option{
		bot.WithMiddlewares(allownedChat(cfg, logger)),
	}

	b, err := bot.New(cfg.Token, opts...)
	if err != nil {
		panic(err)
	}

	userHandler := newUserHandler(logger, userSvc)
	categoryHandler := newCategoryHandler(logger, categorySvc)

	registerHandlers(b, userHandler, categoryHandler)

	b.Start(ctx)

	return nil
}
