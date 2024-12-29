package commandhandler

import (
	"fmt"
	"shikimori-notificator/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Enablecommand(update *tgbotapi.Update, user *models.User, args []string) {
	disabledCommandName := args[1]
	disabledCommandFunction, ok := disabledCommands[disabledCommandName]
	if !ok {
		h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Команда %s не найдена.", disabledCommandName)))
		return
	}
	for _, cmd := range h.commands {
		if cmd.Name == disabledCommandName {
			cmd.Function = disabledCommandFunction
			delete(disabledCommands, disabledCommandName)
			h.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "Команда включена."))
			return
		}
	}
	h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Команда %s не найдена.", disabledCommandName)))
}
