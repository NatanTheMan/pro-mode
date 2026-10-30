package telegram

import (
	"context"

	"github.com/NatanTheMan/pro-mode/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var errMsg map[error]string = map[error]string{
	domain.ErrInvalidName:       "Nome de usuario invalido",
	domain.ErrMaxNameSize:       "Nome excedeu o limite de caracteres: 64",
	domain.ErrMaxUsernameSize:   "Username excedeu o limite de caracteres: 32",
	domain.ErrUserAlreadyExists: "Usuario ja registrado",
}

func reply(b *bot.Bot, ctx context.Context, u *models.Update, msg string) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    u.Message.Chat.ID,
		Text:      msg,
		ParseMode: models.ParseModeMarkdown,

		ReplyParameters: &models.ReplyParameters{
			MessageID: u.Message.ID,
		},
	})
	return err
}
