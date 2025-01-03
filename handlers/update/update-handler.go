package updatehandler

import (
	"fmt"
	"log"
	callbackhandler "shikimori-notificator/handlers/callback"
	commandhandler "shikimori-notificator/handlers/command"
	"shikimori-notificator/models"
	profilenotificator "shikimori-notificator/workers/profile-notificator"
	topicnotificator "shikimori-notificator/workers/topic-notificator"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	"gorm.io/gorm"
)

type UpdateHandler struct {
	Bot   *tgbotapi.BotAPI
	Shiki *shikimori.Client

	Database *gorm.DB

	CommandHandler  *commandhandler.CommandHandler
	CallbackHandler *callbackhandler.CallbackHandler
}

func New(bot *tgbotapi.BotAPI, shiki *shikimori.Client, db *gorm.DB, topicNotificator *topicnotificator.TopicNotificator, profileNotificator *profilenotificator.ProfileNotificator) *UpdateHandler {
	return &UpdateHandler{
		Bot:      bot,
		Shiki:    shiki,
		Database: db,

		CommandHandler:  commandhandler.NewCommandHandler(bot, shiki, topicNotificator, profileNotificator, db),
		CallbackHandler: callbackhandler.NewCallbackHandler(bot, shiki, topicNotificator, profileNotificator, db),
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
			// todo обработка текстового сообщения. Ссылки на топик, чистый айди, ссылки на комментарий и прочее
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
