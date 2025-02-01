package commandhandler

import (
	"fmt"
	"shikimori-notificator/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Topics(update *tgbotapi.Update, user *models.User, args []string) {
	var trackedTopics []models.TrackedTopic
	err := h.Database.Find(&trackedTopics, "user_id = ?", user.ID).Error
	if err != nil {
		panic(err)
	}
	if len(trackedTopics) == 0 {
		h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "У тебя пока нет отслеживаемых топиков."))
		return
	}
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Отслеживаемые топики")
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, trackedTopic := range trackedTopics {
		topic, err := h.ShikiDB.GetTopic(trackedTopic.TopicID)
		if err != nil {
			panic(err)
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(topic.TopicTitle, fmt.Sprintf("topic %d", topic.ID)),
		))
	}
	msg.ReplyMarkup = keyboard
	h.Bot.Send(msg)
}
