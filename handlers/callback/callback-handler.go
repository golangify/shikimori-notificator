package callbackhandler

import (
	"fmt"
	"html"
	"regexp"
	"shikimori-notificator/models"
	topicnotificator "shikimori-notificator/workers/topic-notificator"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	"gorm.io/gorm"
)

type Callback struct {
	Regexp *regexp.Regexp
	Func   func(*Callback, *tgbotapi.Update, *models.User)
}

type CallbackHandler struct {
	Bot              *tgbotapi.BotAPI
	Shiki            *shikimori.Client
	TopicNotificator *topicnotificator.TopicNotificator

	Database *gorm.DB

	callbacks []Callback
}

func NewCallbackHandler(bot *tgbotapi.BotAPI, shiki *shikimori.Client, topicNotificator *topicnotificator.TopicNotificator, db *gorm.DB) *CallbackHandler {
	h := &CallbackHandler{
		Bot:              bot,
		Shiki:            shiki,
		TopicNotificator: topicNotificator,
		Database:         db,
	}

	h.callbacks = []Callback{
		{
			Regexp: regexp.MustCompile(`^add_topic_to_tracking (\d+)$`),
			Func:   h.AddTopicToTracking,
		},
		{
			Regexp: regexp.MustCompile(`^delete_topic_from_tracking (\d+)$`),
			Func:   h.DeleteTopicFromTracking,
		},
		{
			Regexp: regexp.MustCompile(`^topic (\d+)$`),
			Func:   h.Topic,
		},
	}

	return h
}

func (h *CallbackHandler) Process(update *tgbotapi.Update, user *models.User) {
	for _, clb := range h.callbacks {
		if clb.Regexp.MatchString(update.CallbackQuery.Data) {
			clb.Func(&clb, update, user)
			return
		}
	}
	msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Ошибка. Для <code>%s</code> не найдено обработчиков.", html.EscapeString(update.CallbackQuery.Data)))
	msg.ParseMode = tgbotapi.ModeHTML
	h.Bot.Send(msg)
}
