// Package telegram
package telegram

import (
	"context"
	"log/slog"

	"github.com/NatanTheMan/pro-mode/internal/config"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func allownedChat(cfg *config.Config, logger *slog.Logger) func(bot.HandlerFunc) bot.HandlerFunc {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, u *models.Update) {
			msg := u.Message
			if msg != nil && msg.Chat.ID != cfg.ChatID {
				logger.Warn("not allowned chat", "chat_id", msg.Chat.ID)
				return
			}

			next(ctx, b, u)
		}
	}
}
