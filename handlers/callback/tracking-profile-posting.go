package callbackhandler

import (
	"shikimori-notificator/models"
	profileconstructor "shikimori-notificator/view/constructors/profile"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CallbackHandler) setTrackingProfilePosting(с *Callback, update *tgbotapi.Update, user *models.User) {
	profileID, err := strconv.ParseUint(с.Regexp.FindStringSubmatch(update.CallbackQuery.Data)[1], 10, 32)
	if err != nil {
		panic(err)
	}
	profile, err := h.ShikiDB.GetProfile(uint(profileID))
	if err != nil {
		panic(err)
	}
	value, err := strconv.ParseBool(с.Regexp.FindStringSubmatch(update.CallbackQuery.Data)[2])
	if err != nil {
		panic(err)
	}

	var trackedProfile models.TrackedProfile
	if err := h.Database.First(&trackedProfile, "profile_id = ?", profileID).Error; err != nil {
		panic(err)
	}

	err = h.Database.Model(&trackedProfile).UpdateColumn("track_posting", value).Error
	if err != nil {
		panic(err)
	}

	h.Bot.Send(tgbotapi.NewEditMessageReplyMarkup(
		update.FromChat().ID,
		update.CallbackQuery.Message.MessageID,
		*profileconstructor.ToInlineKeyboard(user, &trackedProfile, profile),
	))

	msg := tgbotapi.NewMessage(update.FromChat().ID, "")
	msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
	if value {
		msg.Text = "Теперь постинг пользователя отслеживается."
	} else {
		msg.Text = "Постинг пользователя больше не отслеживается."
	}
	h.Bot.Send(msg)
}
