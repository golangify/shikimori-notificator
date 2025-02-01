package commandhandler

import (
	"fmt"
	"shikimori-notificator/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Profiles(update *tgbotapi.Update, user *models.User, args []string) {
	var trackedProfiles []models.TrackedProfile
	err := h.Database.Find(&trackedProfiles, "user_id = ?", user.ID).Error
	if err != nil {
		panic(err)
	}
	if len(trackedProfiles) == 0 {
		h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "У тебя пока нет отслеживаемых профилей."))
		return
	}
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Отслеживаемые профили")
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, trackedProfile := range trackedProfiles {
		profile, err := h.ShikiDB.GetProfile(trackedProfile.ProfileID)
		if err != nil {
			panic(err)
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(profile.Nickname, fmt.Sprintf("profile %d", profile.ID)),
		))
	}
	msg.ReplyMarkup = keyboard
	h.Bot.Send(msg)
}
