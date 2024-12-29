package callbackhandler

import (
	"fmt"
	"shikimori-notificator/models"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CallbackHandler) AddTopicToTracking(с *Callback, update *tgbotapi.Update, user *models.User) {
	topicID, err := strconv.ParseUint(с.Regexp.FindStringSubmatch(update.CallbackQuery.Data)[1], 10, 32)
	if err != nil {
		panic(err)
	}
	if h.TopicNotificator.IsUserTrackingTopic(user.ID, uint(topicID)) {
		h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Ты уже отслеживаешь этот топик."))
		return
	}
	err = h.TopicNotificator.AddTrackingTopic(user.ID, uint(topicID))
	if err != nil {
		panic(err)
	}
	h.Bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("удалить из отслеживаемого", fmt.Sprintf("delete_topic_from_tracking %d", topicID)),
		),
	)))
	h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Топик добавлен в отслеживаемые."))
}
