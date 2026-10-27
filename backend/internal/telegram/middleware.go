// Package telegram
package telegram

import (
	"context"
	"fmt"

	"github.com/NatanTheMan/pro-mode/internal/config"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func allownedChat(cfg *config.Config) func(bot.HandlerFunc) bot.HandlerFunc {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if update.Message != nil && update.Message.Chat.ID != cfg.ChatID {
				fmt.Println("chat not allowned\n\n\n\n\n\n")
				return
			}

			next(ctx, b, update)
		}
	}
}
