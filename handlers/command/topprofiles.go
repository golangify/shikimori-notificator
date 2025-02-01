package commandhandler

import (
	"fmt"
	"shikimori-notificator/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Topprofiles(update *tgbotapi.Update, user *models.User, args []string) {
	var topTrackedProfiles []models.TrackedProfile
	err := h.Database.Model(&models.TrackedProfile{}).
		Select("*, COUNT(profile_id) as count").
		Group("profile_id").
		Order("count DESC").
		Limit(10).
		Scan(&topTrackedProfiles).Error
	if err != nil {
		panic(err)
	}
	if len(topTrackedProfiles) == 0 {
		panic("сейчас никто ничего не отслеживает")
	}
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Топ 10 Профилей")
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, topTrackedProfile := range topTrackedProfiles {
		profile, err := h.ShikiDB.GetProfile(topTrackedProfile.ProfileID)
		if err != nil {
			panic(err)
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(profile.Nickname, fmt.Sprint("profile ", profile.ID))),
		)
	}
	msg.ReplyMarkup = keyboard
	h.Bot.Send(msg)
}
