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

type UserHandler struct {
	logger  *slog.Logger
	service *service.UserService
}

func newUserHandler(l *slog.Logger, svc *service.UserService) *UserHandler {
	return &UserHandler{logger: l, service: svc}
}

func (h *UserHandler) Register(ctx context.Context, b *bot.Bot, update *models.Update) {
	from := update.Message.From

	user, err := h.service.Register(ctx, from.FirstName, from.Username, from.ID)
	if err != nil {
		if !errors.Is(err, domain.ErrUserAlreadyExists) &&
			!errors.Is(err, domain.ErrInvalidName) {
			h.logger.Error("unspected error", "error", err)
			return
		}

		reply(ctx, b, update, h.logger, errMsg[err])

		return
	}
	msg := fmt.Sprintf("Parabens *%s*, voce esta registrado no ProMode\\!\nFuncao: *%s*\nUsername: @%s",
		user.Name,
		user.Role,
		user.Username)

	h.logger.Info("user registered", "user_id", update.Message.From.ID)

	replyMarkdown(ctx, b, update, h.logger, msg)
}
