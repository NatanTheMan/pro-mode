package telegram

import (
	"context"
	"fmt"

	"github.com/NatanTheMan/pro-mode/internal/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	svc *service.UserService
}

func newHandler(svc *service.UserService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(ctx context.Context, b *bot.Bot, update *models.Update) {
	user, err := h.svc.Register(ctx,
		update.Message.From.FirstName,
		update.Message.From.Username,
		update.Message.From.ID)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})
		return
	}
	msg := fmt.Sprintf("Parabens *%s*, voce esta registrado no ProMode\\!\nFuncao: *%s*\nUsername: *%s*",
		user.Name,
		user.Role,
		user.Username)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      msg,
		ParseMode: models.ParseModeMarkdown,

		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ID,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}
