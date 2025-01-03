package callbackhandler

import (
	"shikimori-notificator/models"
	profileconstructor "shikimori-notificator/view/constructors/profile"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
)

func (h *CallbackHandler) Profile(с *Callback, update *tgbotapi.Update, user *models.User) {
	profileID, err := strconv.ParseUint(с.Regexp.FindStringSubmatch(update.CallbackQuery.Data)[1], 10, 32)
	if err != nil {
		panic(err)
	}
	profile, err := h.Profilenotificator.GetUserProfile(uint(profileID))
	if err != nil {
		if err == shikimori.ErrNotFound {
			msg := tgbotapi.NewMessage(update.FromChat().ID, "Пользователь не найден.")
			h.Bot.Send(msg)
			return
		}
		panic(err)
	}

	msg := tgbotapi.NewMessage(update.FromChat().ID, profileconstructor.ToMessageText(profile))
	msg.ParseMode = tgbotapi.ModeHTML

	keyboard := profileconstructor.ToInlineKeyboard(profile, h.Profilenotificator.IsUserTrackingProfile(user.ID, profile.ID))
	msg.ReplyMarkup = keyboard

	h.Bot.Send(msg)
}
