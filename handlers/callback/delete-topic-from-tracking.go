package callbackhandler

import (
	"fmt"
	"shikimori-notificator/models"
	topicconstructor "shikimori-notificator/view/constructors/topic"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golangify/go-shiki-api/types"
)

func (h *CallbackHandler) DeleteTopicFromTracking(с *Callback, update *tgbotapi.Update, user *models.User) {
	topicID, err := strconv.ParseUint(с.Regexp.FindStringSubmatch(update.CallbackQuery.Data)[1], 10, 32)
	if err != nil {
		panic(err)
	}
	if !h.TopicNotificator.IsUserTrackingTopic(user.ID, uint(topicID)) {
		h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Топика с id %d нет в твоём списке отслеживаемого.", topicID)))
		return
	}
	h.TopicNotificator.DeleteTrackingTopic(user.ID, uint(topicID))
	h.Bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, *topicconstructor.ToInlineKeyboard(&types.Topic{ID: uint(topicID)}, false)))
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Топик удалён из отслеживаемого.")
	msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
	h.Bot.Send(msg)
}
