package telegram

import (
	"context"

	"github.com/NatanTheMan/pro-mode/internal/config"
	"github.com/NatanTheMan/pro-mode/internal/service"
	"github.com/go-telegram/bot"
)

func Run(ctx context.Context, cfg *config.Config, svc *service.UserService) error {
	opts := []bot.Option{
		bot.WithMiddlewares(allownedChat(cfg)),
	}

	b, err := bot.New(cfg.Token, opts...)
	if err != nil {
		panic(err)
	}

	handler := newHandler(svc)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/registrar", bot.MatchTypePrefix, handler.Register)

	b.Start(ctx)

	return nil
}
