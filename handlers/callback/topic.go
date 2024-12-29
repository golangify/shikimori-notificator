package callbackhandler

import (
	"fmt"
	"html"
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
	topic, err := h.TopicNotificator.GetTopic(uint(topicID))
	if err != nil {
		if err == shikimori.ErrNotFound {
			msg := tgbotapi.NewMessage(update.FromChat().ID, "Топик не найден.")
			h.Bot.Send(msg)
			return
		}
		panic(err)
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("<a href='%s'>%s</a>\n\n%s",
		shikimori.ShikiSchema+"://"+shikimori.ShikiDomain+topic.Forum.URL+"/"+strconv.FormatUint(topicID, 10),
		topic.TopicTitle,
		html.EscapeString(topic.Body),
	))
	if len(msg.Text) > 3900 {
		msg.Text = msg.Text[:3900]
	}

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
