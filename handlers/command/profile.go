package commandhandler

import (
	"shikimori-notificator/models"
	profileconstructor "shikimori-notificator/view/constructors/profile"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	shikitypes "github.com/golangify/go-shiki-api/types"
)

func (h *CommandHandler) Profile(update *tgbotapi.Update, user *models.User, args []string) {
	profileID, _ := strconv.ParseUint(args[1], 10, 32)
	var (
		profile *shikitypes.UserProfile
		err     error
	)
	if profileID == 0 {
		profile, err = h.ProfileNotificator.GetUserProfileByNickname(args[1])
	} else {
		profile, err = h.ProfileNotificator.GetUserProfile(uint(profileID))
	}
	if err != nil {
		if err == shikimori.ErrNotFound {
			msg := tgbotapi.NewMessage(update.FromChat().ID, "Пользователь не найден.")
			h.Bot.Send(msg)
			return
		}
		panic(err)
	}

	trackingProfile, err := h.ProfileNotificator.GetTrakingProfile(profile.ID, user.ID)
	if err != nil {
		panic(err)
	}

	msg := tgbotapi.NewMessage(update.FromChat().ID, profileconstructor.ToMessageText(profile))
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = profileconstructor.ToInlineKeyboard(user, trackingProfile, profile)

	h.Bot.Send(msg)
}
