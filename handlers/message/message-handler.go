package messagehandler

import (
	"regexp"
	"shikimori-notificator/models"
	profilenotificator "shikimori-notificator/workers/profile-notificator"
	shikidb "shikimori-notificator/workers/shiki-db"
	topicnotificator "shikimori-notificator/workers/topic-notificator"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	"gorm.io/gorm"
)

type Message struct {
	Name             string
	Description      string
	Usage            string
	Regexp           *regexp.Regexp
	ActivatorRegexps []*regexp.Regexp
	Func             func(update *tgbotapi.Update, user *models.User, args []string)
}

func (m *Message) Help() string {
	helpText := "Используй: " + m.Usage + " - " + m.Description
	if m.Func == nil {
		helpText += " (действие временно недоступно)"
	}
	return helpText
}

type MessageHandler struct {
	Bot      *tgbotapi.BotAPI
	Shiki    *shikimori.Client
	ShikiDB  *shikidb.ShikiDB
	Database *gorm.DB

	TopicNotificator   *topicnotificator.TopicNotificator
	Profilenotificator *profilenotificator.ProfileNotificator

	messages []Message
}

func NewMessageHandler(bot *tgbotapi.BotAPI, shiki *shikimori.Client, shikidb *shikidb.ShikiDB, database *gorm.DB, topicNotificator *topicnotificator.TopicNotificator, profileNotificator *profilenotificator.ProfileNotificator) *MessageHandler {
	h := &MessageHandler{
		Bot:      bot,
		Shiki:    shiki,
		ShikiDB:  shikidb,
		Database: database,

		TopicNotificator:   topicNotificator,
		Profilenotificator: profileNotificator,
	}

	h.messages = []Message{
		{
			Name:        "from-comment",
			Description: "получить топик/профиль по ссылке комментария",
			Usage:       "https://shikimori.one/comments/10478513",
			ActivatorRegexps: []*regexp.Regexp{
				regexp.MustCompile(`comments/(\d+)`),
			},
			Regexp: regexp.MustCompile("^" + shikimori.ShikiSchema + `:\/\/` + shikimori.ShikiDomain + `/comments/(\d+)$`),
			Func:   h.FromComment,
		},
	}

	return h
}

func (h *MessageHandler) Process(update *tgbotapi.Update, user *models.User) {
	if update.Message.Text == "" {
		h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Понимаю только текстовые сообщения."))
		return
	}

	for _, msg := range h.messages {
		if msg.Regexp.MatchString(update.Message.Text) {
			if msg.Func == nil {
				h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Действие временно не доступно."))
				return
			}
			msg.Func(update, user, msg.Regexp.FindAllStringSubmatch(update.Message.Text, -1)[0])
			return
		}
		for _, activatorRegexp := range msg.ActivatorRegexps {
			if activatorRegexp.MatchString(update.Message.Text) {
				h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, msg.Help()))
				return
			}
		}
	}
	h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Не удалось распознать действие."))
}
