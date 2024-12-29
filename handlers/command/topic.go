package commandhandler

import (
	"fmt"
	"html"
	"shikimori-notificator/models"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
)

func (h *CommandHandler) Topic(update *tgbotapi.Update, user *models.User, args []string) {
	topicID, _ := strconv.ParseUint(args[1], 10, 32)
	topic, err := h.TopicNotificator.GetTopic(uint(topicID))
	if err != nil {
		if err == shikimori.ErrNotFound {
			msg := tgbotapi.NewMessage(update.FromChat().ID, "Топик не найден.")
			h.Bot.Send(msg)
			return
		}
		panic(err)
	}

	messageText := fmt.Sprintf("<a href='%s'>%s</a>\n\n%s",
		shikimori.ShikiSchema+"://"+shikimori.ShikiDomain+topic.Forum.URL+"/"+args[1],
		topic.TopicTitle,
		html.EscapeString(topic.Body),
	)
	if len(messageText) > 3900 {
		messageText = messageText[:3900]
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
	if h.TopicNotificator.IsUserTrackingTopic(user.ID, topic.ID) {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("удалить из отслеживаемого", fmt.Sprintf("delete_topic_from_tracking %d", topicID)),
			),
		)
	} else {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("добавить в отслеживаемое", fmt.Sprintf("add_topic_to_tracking %d", topicID)),
			),
		)
	}
	msg.ParseMode = tgbotapi.ModeHTML

	h.Bot.Send(msg)
}
