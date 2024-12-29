package commandhandler

import (
	"fmt"
	"shikimori-notificator/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	disabledCommands map[string]commandFunction
)

func (h *CommandHandler) Disablecommand(update *tgbotapi.Update, user *models.User, args []string) {
	commandNameToDisable := args[1]
	if commandNameToDisable == "enablecommand" || commandNameToDisable == "disablecommand" {
		h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Очень смешно."))
		return
	}
	for _, cmd := range h.commands {
		if cmd.Name == commandNameToDisable {
			if cmd.Function == nil {
				h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Команда уже отключена."))
				return
			}
			if disabledCommands == nil {
				disabledCommands = make(map[string]commandFunction)
			}
			disabledCommands[commandNameToDisable] = cmd.Function
			cmd.Function = nil
			h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Команда %s отключена.", cmd.Name)))
			return
		}
	}
	h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Команда с именем %s не найдена.", commandNameToDisable)))
}
