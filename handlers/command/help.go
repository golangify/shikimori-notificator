package commandhandler

import (
	"shikimori-notificator/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Help(update *tgbotapi.Update, user *models.User, args []string) {
	msg := tgbotapi.NewMessage(update.FromChat().ID, "<b>Помощь</b>\n")
	for _, cmd := range h.commands {
		if cmd.Level > user.Level {
			continue // пропускаем команду, т.к. у пользователя недостаточно прав на её использование
		}
		msg.Text += cmd.Help() + "\n\n"
	}
	msg.Text += "А также можно отправить ссылку на комментарий(https://shikimori.one/comments/10478513 пример), бот отправит команду на получение топика/профиля из комментария для отслеживания."
	msg.ParseMode = tgbotapi.ModeHTML
	h.Bot.Send(msg)
}
