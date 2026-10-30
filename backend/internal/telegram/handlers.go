package telegram

import (
	"github.com/go-telegram/bot"
)

func registerHandlers(
	b *bot.Bot,
	userHandler *UserHandler,
	categoryHandler *CategoryHandler,
) {
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/registrar",
		bot.MatchTypePrefix,
		userHandler.Register,
	)

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/add_categoria",
		bot.MatchTypePrefix,
		categoryHandler.CreateCategory,
	)
}
