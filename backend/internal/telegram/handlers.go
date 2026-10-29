package telegram

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/NatanTheMan/pro-mode/internal/domain"
	"github.com/NatanTheMan/pro-mode/internal/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	svc    *service.UserService
	logger *slog.Logger
}

func newHandler(svc *service.UserService, logger *slog.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

func (h *Handler) Register(ctx context.Context, b *bot.Bot, update *models.Update) {
	user, err := h.svc.Register(ctx,
		update.Message.From.FirstName,
		update.Message.From.Username,
		update.Message.From.ID)
	if err != nil {
		if !errors.Is(err, domain.ErrUserAlreadyExists) &&
			!errors.Is(err, domain.ErrInvalidName) {
			h.logger.Error("unspected error", "error", err)
			return
		}

		var msgErr string
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			msgErr = "Usuario ja registrado"
		} else if errors.Is(err, domain.ErrInvalidName) {
			msgErr = "Nome de usuario invalido"
		} else if errors.Is(err, domain.ErrMaxNameSize) {
			msgErr = "Nome excedeu o limite de caracteres: 64"
		} else if errors.Is(err, domain.ErrMaxUsernameSize) {
			msgErr = "Username excedeu o limite de caracteres: 32"
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   msgErr,
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})
		if err != nil {
			h.logger.Error("failed to send message", "error", err)
		}

		return
	}
	msg := fmt.Sprintf("Parabens *%s*, voce esta registrado no ProMode\\!\nFuncao: *%s*\nUsername: *%s*",
		user.Name,
		user.Role,
		user.Username)

	h.logger.Info("user registered", "user_id", update.Message.From.ID)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      msg,
		ParseMode: models.ParseModeMarkdown,

		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ID,
		},
	})
	if err != nil {
		h.logger.Error("failed to send message", "error", err)
	}
}
