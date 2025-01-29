package profileconstructor

import (
	"fmt"
	"html"
	"shikimori-notificator/models"

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

func ToInlineKeyboard(user *models.User, tracked *models.TrackedProfile, targetProfile *shikitypes.UserProfile) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	if tracked != nil {
		if tracked.TrackPosting {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("прекратить отслеживать постинг", fmt.Sprint("set_tracking_profile_posting ", tracked.ProfileID, " false")),
			))
		} else {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("отслеживать постинг", fmt.Sprint("set_tracking_profile_posting ", tracked.ProfileID, " true")),
			))
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("прекратить отслеживать профиль", fmt.Sprint("delete_profile_from_tracking ", tracked.ProfileID)),
		))
	} else {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("отслеживать профиль", fmt.Sprint("add_profile_to_tracking ", targetProfile.ID)),
		))
	}
	return &keyboard
}
