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
	svc *service.UserService,
	logger *slog.Logger,
) error {
	opts := []bot.Option{
		bot.WithMiddlewares(allownedChat(cfg, logger)),
	}

	b, err := bot.New(cfg.Token, opts...)
	if err != nil {
		panic(err)
	}

	handler := newHandler(svc, logger)

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/registrar",
		bot.MatchTypePrefix,
		handler.Register,
	)

	b.Start(ctx)

	return nil
}
