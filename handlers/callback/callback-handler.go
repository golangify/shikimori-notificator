package callbackhandler

import (
	"fmt"
	"html"
	"regexp"
	"shikimori-notificator/models"
	profilenotificator "shikimori-notificator/workers/profile-notificator"
	shikidb "shikimori-notificator/workers/shiki-db"
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
	Bot      *tgbotapi.BotAPI
	Shiki    *shikimori.Client
	ShikiDB  *shikidb.ShikiDB
	Database *gorm.DB

	TopicNotificator   *topicnotificator.TopicNotificator
	ProfileNotificator *profilenotificator.ProfileNotificator

	callbacks []Callback
}

func NewCallbackHandler(bot *tgbotapi.BotAPI, shiki *shikimori.Client, shikidb *shikidb.ShikiDB, db *gorm.DB, topicNotificator *topicnotificator.TopicNotificator, profileNotificator *profilenotificator.ProfileNotificator) *CallbackHandler {
	h := &CallbackHandler{
		Bot:      bot,
		Shiki:    shiki,
		ShikiDB:  shikidb,
		Database: db,

		TopicNotificator:   topicNotificator,
		ProfileNotificator: profileNotificator,
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
		{
			Regexp: regexp.MustCompile(`^add_profile_to_tracking (\d+)$`),
			Func:   h.AddProfileToTracking,
		},
		{
			Regexp: regexp.MustCompile(`^delete_profile_from_tracking (\d+)$`),
			Func:   h.DeleteProfileFromTracking,
		},
		{
			Regexp: regexp.MustCompile(`^profile (\d+)$`),
			Func:   h.Profile,
		},
		{
			Regexp: regexp.MustCompile(`set_tracking_profile_posting (\d+) (true|false)`),
			Func:   h.setTrackingProfilePosting,
		},
	}

	return h
}

func (h *CallbackHandler) Process(update *tgbotapi.Update, user *models.User) {
	go h.Bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "..."))
	for _, clb := range h.callbacks {
		if clb.Regexp.MatchString(update.CallbackQuery.Data) {
			if clb.Func == nil {
				h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Кнопка временно недоступна."))
				return
			}
			clb.Func(&clb, update, user)
			return
		}
	}
	msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Ошибка. Для <code>%s</code> не найдено обработчиков.", html.EscapeString(update.CallbackQuery.Data)))
	msg.ParseMode = tgbotapi.ModeHTML
	h.Bot.Send(msg)
}
