package commandhandler

import (
	"fmt"
	"shikimori-notificator/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Clearcache(update *tgbotapi.Update, user *models.User, args []string) {
	numDeleted := h.ShikiDB.ClearCache()

	h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Кэш очищен. Объектов удалено: %d", numDeleted)))
}
