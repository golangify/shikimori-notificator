package profilenotificator

import (
	"log"
	"shikimori-notificator/models"
	commentconstructor "shikimori-notificator/view/constructors/comment"
	"slices"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikitypes "github.com/golangify/go-shiki-api/types"
)

func (n *ProfileNotificator) notifyProfiles() error {
	var trackedProfiles []models.TrackedProfile
	err := n.Database.Find(&trackedProfiles).Order("last_comment_id").Distinct("profile_id").Error
	if err != nil {
		log.Println(err)
		return err
	}
	t := time.NewTicker(time.Second * 2)
	defer t.Stop()
	for _, trackedProfile := range trackedProfiles {
		<-t.C
		userProfile, err := n.GetUserProfile(trackedProfile.ProfileID)
		if err != nil {
			log.Println(err)
			continue
		}

		lastComments, err := n.Shiki.GetComments(userProfile.ID, shikitypes.TypeUser, 1, 20, true)
		if err != nil {
			log.Println(err)
			continue
		}
		var newComments []shikitypes.Comment
		for _, comment := range lastComments {
			if comment.ID > trackedProfile.LastCommentID {
				newComments = append(newComments, comment)
			}
		}
		if len(newComments) == 0 {
			continue
		}
		slices.Reverse(newComments)

		var usersTrackedProfile []models.TrackedProfile
		err = n.Database.Preload("User").Find(&usersTrackedProfile, "profile_id = ?", userProfile.ID).Error
		if err != nil {
			log.Println(err)
			continue
		}
		for _, newComment := range newComments {
			msg := tgbotapi.NewMessage(0, commentconstructor.ProfileToMessageText(&newComment, userProfile))
			msg.ParseMode = tgbotapi.ModeHTML
			msg.DisableWebPagePreview = true
			for _, userTrackedProfile := range usersTrackedProfile {
				if !n.filter.Ok(newComment.ID, userTrackedProfile.User.ID) {
					continue
				}
				msg.BaseChat.ChatID = userTrackedProfile.User.TgID
				_, err := n.Bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
		err = n.Database.Model(&models.TrackedProfile{}).Where("profile_id = ?", userProfile.ID).UpdateColumn("last_comment_id", newComments[len(newComments)-1].ID).Error
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
