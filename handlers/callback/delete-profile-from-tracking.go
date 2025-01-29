package callbackhandler

import (
	"fmt"
	"shikimori-notificator/models"
	profileconstructor "shikimori-notificator/view/constructors/profile"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CallbackHandler) DeleteProfileFromTracking(с *Callback, update *tgbotapi.Update, user *models.User) {
	profileID, err := strconv.ParseUint(с.Regexp.FindStringSubmatch(update.CallbackQuery.Data)[1], 10, 32)
	if err != nil {
		panic(err)
	}
	profile, err := h.ProfileNotificator.GetUserProfile(uint(profileID))
	if err != nil {
		panic(err)
	}
	if !h.ProfileNotificator.IsUserTrackingProfile(user.ID, uint(profileID)) {
		h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Профиля с id %d нет в твоём списке.", profileID)))
		return
	}
	h.ProfileNotificator.DeleteTrackingProfile(uint(profileID), user.ID)

	trackingProfile, err := h.ProfileNotificator.GetTrackingProfile(profile.ID, user.ID)
	if err != nil {
		panic(err)
	}

	h.Bot.Send(tgbotapi.NewEditMessageReplyMarkup(
		update.FromChat().ID,
		update.CallbackQuery.Message.MessageID,
		*profileconstructor.ToInlineKeyboard(user, trackingProfile, profile),
	))
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Профиль удалён из отслеживаемого.")
	msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
	h.Bot.Send(msg)
}
