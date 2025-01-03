package commandhandler

import (
	"fmt"
	"shikimori-notificator/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Toptopics(update *tgbotapi.Update, user *models.User, args []string) {
	var topTrackedTopics []models.TrackedTopic
	err := h.Database.Model(&models.TrackedTopic{}).
		Select("*, COUNT(topic_id) as count").
		Group("topic_id").
		Order("count DESC").
		Limit(10).
		Scan(&topTrackedTopics).Error
	if err != nil {
		panic(err)
	}
	if len(topTrackedTopics) == 0 {
		panic("сейчас никто ничего не отслеживает")
	}
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Топ 10 топиков")
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, topTrackedTopic := range topTrackedTopics {
		topic, err := h.TopicNotificator.GetTopic(topTrackedTopic.TopicID)
		if err != nil {
			panic(err)
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(topic.TopicTitle, fmt.Sprintf("topic %d", topTrackedTopic.TopicID)),
		))
	}
	msg.ReplyMarkup = keyboard
	h.Bot.Send(msg)
}
