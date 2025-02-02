package commandhandler

import (
	"shikimori-notificator/models"
	topicconstructor "shikimori-notificator/view/constructors/topic"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
)

func (h *CommandHandler) Topic(update *tgbotapi.Update, user *models.User, args []string) {
	topicID, _ := strconv.ParseUint(args[2], 10, 32)
	topic, err := h.ShikiDB.GetTopic(uint(topicID))
	if err != nil {
		if err == shikimori.ErrNotFound {
			msg := tgbotapi.NewMessage(update.FromChat().ID, "Топик не найден.")
			h.Bot.Send(msg)
			return
		}
		panic(err)
	}

	messageText := topicconstructor.ToMessageText(topic)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
	msg.ReplyMarkup = topicconstructor.ToInlineKeyboard(topic, h.TopicNotificator.IsUserTrackingTopic(user.ID, topic.ID))
	msg.ParseMode = tgbotapi.ModeHTML

	_, err = h.Bot.Send(msg)
	if err != nil {
		panic(err)
	}
}
