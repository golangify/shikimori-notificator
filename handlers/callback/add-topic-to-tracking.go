package callbackhandler

import (
	"shikimori-notificator/models"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikitypes "github.com/golangify/go-shiki-api/types"
)

func (h *CallbackHandler) AddTopicToTracking(с *Callback, update *tgbotapi.Update, user *models.User) {
	topicID, err := strconv.ParseUint(с.Regexp.FindStringSubmatch(update.CallbackQuery.Data)[1], 10, 32)
	if err != nil {
		panic(err)
	}
	if h.TopicNotificator.IsUserTrackingTopic(user.ID, uint(topicID)) {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Ты уже отслеживаешь этот топик.")
		msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
		h.Bot.Send(msg)
		return
	}
	err = h.TopicNotificator.AddTrackingTopic(user.ID, uint(topicID))
	if err != nil {
		panic(err)
	}
	h.Bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, *h.TopicConstructor.InlineKeyboard(&shikitypes.Topic{ID: uint(topicID)}, true)))
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Топик добавлен в отслеживаемые.")
	msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
	h.Bot.Send(msg)
}
