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
	from := update.Message.From

	user, err := h.svc.Register(ctx, from.FirstName, from.Username, from.ID)
	if err != nil {
		if !errors.Is(err, domain.ErrUserAlreadyExists) &&
			!errors.Is(err, domain.ErrInvalidName) {
			h.logger.Error("unspected error", "error", err)
			return
		}

		err = reply(b, ctx, update, errMsg[err])
		if err != nil {
			h.logger.Error("failed to send message", "error", err)
		}

		return
	}
	msg := fmt.Sprintf("Parabens *%s*, voce esta registrado no ProMode\\!\nFuncao: *%s*\nUsername: @%s",
		user.Name,
		user.Role,
		user.Username)

	h.logger.Info("user registered", "user_id", update.Message.From.ID)

	err = reply(b, ctx, update, msg)
	if err != nil {
		h.logger.Error("failed to send message", "error", err)
	}
}
