package telegram

import (
	"context"
	"log/slog"
	"strings"

	"github.com/NatanTheMan/pro-mode/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var errMsg map[error]string = map[error]string{
	domain.ErrInvalidName:           "Nome de usuario invalido",
	domain.ErrMaxNameSize:           "Nome excedeu o limite de caracteres: 64",
	domain.ErrMaxUsernameSize:       "Username excedeu o limite de caracteres: 32",
	domain.ErrUserAlreadyExists:     "Usuario ja registrado",
	domain.ErrCategoryAlreadyExists: "Esta categoria ja esta cadastrada",
	domain.ErrMaxCategoryNameSize:   "Nome da categoria nao deve exceder 30 caracteres",
	domain.ErrMinCategoryNameSize:   "Nome da categoria deve ter ao menos 3 caracteres",
}

func reply(ctx context.Context, b *bot.Bot, u *models.Update, l *slog.Logger, msg string) {
	baseReply(ctx, b, l, &bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   msg,

		ReplyParameters: &models.ReplyParameters{
			MessageID: u.Message.ID,
		},
	})
}

func replyMarkdown(ctx context.Context, b *bot.Bot, u *models.Update, l *slog.Logger, msg string) {
	baseReply(ctx, b, l, &bot.SendMessageParams{
		ChatID:    u.Message.Chat.ID,
		Text:      msg,
		ParseMode: models.ParseModeMarkdown,

		ReplyParameters: &models.ReplyParameters{
			MessageID: u.Message.ID,
		},
	})
}

func baseReply(ctx context.Context, b *bot.Bot, l *slog.Logger, params *bot.SendMessageParams) {
	_, err := b.SendMessage(ctx, params)
	if err != nil {
		l.Error("message not sended", "error", err)
	}
}

func commandArgs(text string) string {
	_, args, _ := strings.Cut(text, " ")
	return strings.TrimSpace(args)
}
