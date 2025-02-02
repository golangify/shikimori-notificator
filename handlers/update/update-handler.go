package updatehandler

import (
	"fmt"
	"log"
	callbackhandler "shikimori-notificator/handlers/callback"
	commandhandler "shikimori-notificator/handlers/command"
	messagehandler "shikimori-notificator/handlers/message"
	"shikimori-notificator/models"
	"shikimori-notificator/view/constructors/bbcode"
	topicconstructor "shikimori-notificator/view/constructors/topic"
	profilenotificator "shikimori-notificator/workers/profile-notificator"
	shikidb "shikimori-notificator/workers/shiki-db"
	topicnotificator "shikimori-notificator/workers/topic-notificator"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	"gorm.io/gorm"
)

type UpdateHandler struct {
	Bot      *tgbotapi.BotAPI
	Shiki    *shikimori.Client
	ShikiDB  *shikidb.ShikiDB
	Database *gorm.DB

	TopicConstructor *topicconstructor.TopicConstructor

	CommandHandler  *commandhandler.CommandHandler
	CallbackHandler *callbackhandler.CallbackHandler
	MessageHandler  *messagehandler.MessageHandler

	BBCodeParser *bbcode.BBCodeParser
}

func New(bot *tgbotapi.BotAPI, shiki *shikimori.Client, db *gorm.DB, shikidb *shikidb.ShikiDB, topicNotificator *topicnotificator.TopicNotificator, profileNotificator *profilenotificator.ProfileNotificator, topicConstructor *topicconstructor.TopicConstructor) *UpdateHandler {
	return &UpdateHandler{
		Bot:      bot,
		Shiki:    shiki,
		ShikiDB:  shikidb,
		Database: db,

		TopicConstructor: topicConstructor,

		CommandHandler:  commandhandler.NewCommandHandler(bot, shiki, shikidb, db, topicNotificator, profileNotificator, topicConstructor),
		CallbackHandler: callbackhandler.NewCallbackHandler(bot, shiki, shikidb, db, topicNotificator, profileNotificator, topicConstructor),
		MessageHandler:  messagehandler.NewMessageHandler(bot, shiki, shikidb, db, topicNotificator, profileNotificator),
	}
}

func (h *UpdateHandler) Process(update *tgbotapi.Update) {
	defer func() {
		if h.Bot.Debug {
			return
		}
		if err := recover(); err != nil {
			log.Println(err)
			h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("произошла ошибка: %v", err)))
		}
	}()

	user, ok, err := h.validateUserActivity(update)
	if !ok {
		return
	}
	if err != nil {
		panic(err)
	}

	if update.Message != nil {
		if update.Message.IsCommand() {
			h.CommandHandler.Process(update, user)
		} else {
			h.MessageHandler.Process(update, user)
		}
	} else if update.CallbackQuery != nil {
		h.CallbackHandler.Process(update, user)
	} else {
		h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Действие не поддерживается."))
	}
}

func (h *UpdateHandler) validateUserActivity(update *tgbotapi.Update) (*models.User, bool, error) {
	from := update.SentFrom()
	if from == nil {
		return nil, false, nil
	}
	var user models.User
	err := h.Database.First(&user, "tg_id = ?", from.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user.TgID = from.ID
			h.Database.Create(&user)
			return &user, true, nil
		}
		return nil, false, err
	}
	return &user, true, nil
}
