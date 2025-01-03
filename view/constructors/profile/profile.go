package profileconstructor

import (
	"fmt"
	"html"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikitypes "github.com/golangify/go-shiki-api/types"
)

func ToMessageText(p *shikitypes.UserProfile) string {
	return fmt.Sprintf(
		"Профиль <a href='%s'>%s</a>\n\t<i>%s</i>",
		p.URL, html.EscapeString(p.Nickname),
		p.LastOnline,
	)
}

func ToInlineKeyboard(p *shikitypes.UserProfile, isProfileTracking bool) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	if isProfileTracking {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("прекратить отслеживать профиль", fmt.Sprint("delete_profile_from_tracking ", p.ID)),
		))
	} else {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("отслеживать профиль", fmt.Sprint("add_profile_to_tracking ", p.ID)),
		))
	}
	return &keyboard
}
