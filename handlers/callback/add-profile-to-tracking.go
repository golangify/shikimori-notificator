package callbackhandler

import (
	"shikimori-notificator/models"
	profileconstructor "shikimori-notificator/view/constructors/profile"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikitypes "github.com/golangify/go-shiki-api/types"
)

func (h *CallbackHandler) AddProfileToTracking(с *Callback, update *tgbotapi.Update, user *models.User) {
	profileID, err := strconv.ParseUint(с.Regexp.FindStringSubmatch(update.CallbackQuery.Data)[1], 10, 32)
	if err != nil {
		panic(err)
	}
	if h.Profilenotificator.IsUserTrackingProfile(user.ID, uint(profileID)) {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Профиль уже в отслеживаемом!")
		msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
		h.Bot.Send(msg)
		return
	}
	err = h.Profilenotificator.AddTrackingProfile(user.ID, uint(profileID))
	if err != nil {
		panic(err)
	}
	h.Bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, *profileconstructor.ToInlineKeyboard(&shikitypes.UserProfile{ID: uint(profileID)}, true)))
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Профиль добавлен в отслеживаемое.")
	msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
	h.Bot.Send(msg)
}
