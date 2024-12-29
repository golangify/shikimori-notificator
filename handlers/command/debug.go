package commandhandler

import (
	"shikimori-notificator/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Debug(update *tgbotapi.Update, user *models.User, args []string) {
	h.Bot.Debug = !h.Bot.Debug
	if h.Bot.Debug {
		h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Debug включен"))
	} else {
		h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Debug выключен"))
	}
}
