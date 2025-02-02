package callbackhandler

import (
	"shikimori-notificator/models"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
)

func (h *CallbackHandler) Topic(с *Callback, update *tgbotapi.Update, user *models.User) {
	topicID, err := strconv.ParseUint(с.Regexp.FindStringSubmatch(update.CallbackQuery.Data)[1], 10, 32)
	if err != nil {
		panic(err)
	}
	topic, err := h.ShikiDB.GetTopic(uint(topicID))
	if err != nil {
		if err == shikimori.ErrNotFound {
			msg := tgbotapi.NewMessage(update.FromChat().ID, "Топик не найден.")
			h.Bot.Send(msg)
			return
		}
		panic(err)
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, h.TopicConstructor.Text(topic))
	msg.ReplyMarkup = h.TopicConstructor.InlineKeyboard(topic, h.TopicNotificator.IsUserTrackingTopic(user.ID, topic.ID))
	msg.ParseMode = tgbotapi.ModeHTML

	h.Bot.Send(msg)
}
