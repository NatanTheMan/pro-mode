package telegram

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/NatanTheMan/pro-mode/internal/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type CategoryHandler struct {
	logger  *slog.Logger
	service *service.CategoryService
}

func newCategoryHandler(l *slog.Logger, svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{logger: l, service: svc}
}

func (h *CategoryHandler) CreateCategory(ctx context.Context, b *bot.Bot, u *models.Update) {
	name := commandArgs(u.Message.Text)
	if name == "" {
		reply(ctx, b, u, h.logger, "Categoria precisa de um nome. Exemplo:\n/add_categoria Treino")
		return
	}

	category, err := h.service.Create(ctx, name, u.Message.From.ID)
	msg, ok := errMsg[err]
	if err != nil {
		if ok {
			reply(ctx, b, u, h.logger, msg)
			return
		}
		h.logger.Error("unexpected error occured", "error", err)
		reply(ctx, b, u, h.logger, "Um erro inesperado aconteceu")
		return
	}

	msg = fmt.Sprintf("Categoria criada: %s", category.Name)

	reply(ctx, b, u, h.logger, msg)
}
