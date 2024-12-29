package commandhandler

import (
	"fmt"
	"html"
	"shikimori-notificator/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Start(update *tgbotapi.Update, user *models.User, args []string) {
	msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf(
		"<b>Привет, %s</b>!\n\n"+
			"В этом боте можно отслеживать все новые комментарии под темами с сайта shikimori.one.\n\n"+
			"Сводка команд - /help",
		html.EscapeString(update.SentFrom().FirstName),
	))
	msg.ParseMode = tgbotapi.ModeHTML
	h.Bot.Send(msg)
}
