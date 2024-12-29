package commandhandler

import (
	"fmt"
	"html"
	"shikimori-notificator/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Start(update *tgbotapi.Update, user *models.User, args []string) {
	var totalUsersCount, totalTrackedTopicsCount int64
	h.Database.Model(&models.User{}).Count(&totalUsersCount)
	h.Database.Model(&models.TrackedTopic{}).Count(&totalTrackedTopicsCount)

	msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf(
		"<b>Привет, %s</b>!\n\n"+
			"В этом боте можно отслеживать все новые комментарии под темами с сайта shikimori.one.\n\n"+
			"Сводка команд - /help\n\n"+
			"Всего пользователей: %d\n"+
			"Всего отслеживаемых топиков: %d",
		html.EscapeString(update.SentFrom().FirstName),
		totalUsersCount,
		totalTrackedTopicsCount,
	))
	msg.ParseMode = tgbotapi.ModeHTML
	h.Bot.Send(msg)
}
