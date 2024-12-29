package commandhandler

import (
	"fmt"
	"shikimori-notificator/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golangify/go-shiki-api/types"
)

func (h *CommandHandler) Clearcache(update *tgbotapi.Update, user *models.User, args []string) {
	numCached := len(h.TopicNotificator.CachedTopics)
	h.TopicNotificator.Mu.Lock()
	h.TopicNotificator.CachedTopics = make(map[uint]*types.Topic)
	h.TopicNotificator.Mu.Unlock()
	h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Кэш очищен. Объектов удалено: %d", numCached)))
}
