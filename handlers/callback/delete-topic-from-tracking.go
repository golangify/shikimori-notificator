package callbackhandler

import (
	"fmt"
	"shikimori-notificator/models"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	h.Database.Where("topic_id = ? AND user_id = ?", topicID, user.ID).Delete(&models.TrackedTopic{})
	h.Bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("добавить в отслеживаемое", fmt.Sprintf("add_topic_to_tracking %d", topicID)),
		),
	)))
	h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Топик удалён из отслеживаемого."))
}
