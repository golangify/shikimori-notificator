package callbackhandler

import (
	"fmt"
	"shikimori-notificator/models"
	profileconstructor "shikimori-notificator/view/constructors/profile"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikitypes "github.com/golangify/go-shiki-api/types"
)

func (h *CallbackHandler) DeleteProfileFromTracking(с *Callback, update *tgbotapi.Update, user *models.User) {
	profileID, err := strconv.ParseUint(с.Regexp.FindStringSubmatch(update.CallbackQuery.Data)[1], 10, 32)
	if err != nil {
		panic(err)
	}
	if !h.Profilenotificator.IsUserTrackingProfile(user.ID, uint(profileID)) {
		h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Профиля с id %d нет в твоём списке.", profileID)))
		return
	}
	h.Profilenotificator.DeleteTrackingProfile(uint(profileID), user.ID)
	h.Bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, *profileconstructor.ToInlineKeyboard(&shikitypes.UserProfile{ID: uint(profileID)}, false)))
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Профиль удалён из отслеживаемого.")
	msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
	h.Bot.Send(msg)
}
