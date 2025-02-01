package profilenotificator

import (
	"log"
	"shikimori-notificator/models"
	commentconstructor "shikimori-notificator/view/constructors/comment"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	shikitypes "github.com/golangify/go-shiki-api/types"
)

func (n *ProfileNotificator) notifyPosting() error {
	t := time.NewTicker(time.Second * 2)
	defer t.Stop()
	trackedProfiles := make([]models.TrackedProfile, 0)
	if err := n.Database.Find(&trackedProfiles, "track_posting = ?", true).Order("last_posted_comment_id").Distinct("profile_id").Error; err != nil {
		log.Println(err)
		return err
	}
	for _, trackedProfile := range trackedProfiles {
		userProfile, err := n.GetUserProfile(trackedProfile.ProfileID)
		if err != nil {
			log.Println(err)
			continue
		}

		var usersTrackedProfile []models.TrackedProfile
		err = n.Database.Preload("User").Find(&usersTrackedProfile, "profile_id = ? and track_posting = ?", userProfile.ID, true).Error
		if err != nil {
			log.Println(err)
			continue
		}

		commentIDs, err := n.GetLast20PostedCommentIDs(userProfile.Nickname)
		if err != nil {
			log.Println(err)
			continue
		}
		var newCommentIDs []uint
		for _, commentID := range commentIDs {
			if commentID > trackedProfile.LastPostedCommentID {
				newCommentIDs = append(newCommentIDs, commentID)
			}
		}
		if len(newCommentIDs) == 0 {
			continue
		}
		var newComments []shikitypes.Comment
		for _, newCommentID := range newCommentIDs {
			<-t.C
			commment, err := n.GetComment(newCommentID)
			if err != nil {
				if err == shikimori.ErrNotFound {
					continue
				}
				log.Println(err)
				continue
			}
			newComments = append(newComments, *commment)
		}

		for _, newComment := range newComments {
			msg := tgbotapi.NewMessage(0, "")
			if newComment.CommentableType == shikitypes.TypeUser {
				msg.Text = commentconstructor.ProfileToMessageText(&newComment, userProfile)
			} else {
				topic, err := n.TopicNotificator.GetTopic(newComment.CommentableID)
				if err != nil {
					log.Println(err)
					continue
				}
				msg.Text = commentconstructor.TopicToMessageText(&newComment, topic)
			}
			msg.ParseMode = tgbotapi.ModeHTML
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
		err = n.Database.Model(&models.TrackedProfile{}).Where("profile_id = ?", userProfile.ID).UpdateColumn("last_posted_comment_id", newComments[len(newComments)-1].ID).Error
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
