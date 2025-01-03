package profilenotificator

import (
	"fmt"
	"html"
	"log"
	"shikimori-notificator/models"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	shikitypes "github.com/golangify/go-shiki-api/types"
	"gorm.io/gorm"
)

type ProfileNotificator struct {
	Shiki    *shikimori.Client
	Bot      *tgbotapi.BotAPI
	Database *gorm.DB

	CachedProfiles           map[uint]*shikitypes.UserProfile
	CachedProfilesByNickname map[string]*shikitypes.UserProfile
	Mu                       sync.Mutex
	ticker                   *time.Ticker
}

func NewProfileNotificator(shiki *shikimori.Client, bot *tgbotapi.BotAPI, database *gorm.DB) *ProfileNotificator {
	n := &ProfileNotificator{
		Shiki:    shiki,
		Bot:      bot,
		Database: database,

		CachedProfiles:           make(map[uint]*shikitypes.UserProfile),
		CachedProfilesByNickname: make(map[string]*shikitypes.UserProfile),
	}

	return n
}

func (n *ProfileNotificator) Run() {
	n.ticker = time.NewTicker(time.Minute)
	for range n.ticker.C {
		var trackedProfiles []models.TrackedProfile
		err := n.Database.Find(&trackedProfiles).Order("last_comment_id").Distinct("profile_id").Error
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(trackedProfiles)
		t := time.NewTicker(time.Second * 2)
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
			for i, j := 0, len(newComments)-1; i < j; i, j = i+1, j-1 {
				newComments[i], newComments[j] = newComments[j], newComments[i]
			}

			var usersTrackedProfile []models.TrackedProfile
			err = n.Database.Preload("User").Find(&usersTrackedProfile, "profile_id = ?", userProfile.ID).Error
			if err != nil {
				log.Println(err)
				continue
			}
			for _, newComment := range newComments {
				msg := tgbotapi.NewMessage(0, fmt.Sprintf(
					"<a href='%s'>%s</a> в профиле <a href='%s'>%s</a>\n\n%s",
					newComment.User.URL, html.EscapeString(newComment.User.Nickname),
					userProfile.URL, html.EscapeString(userProfile.Nickname),
					html.EscapeString(newComment.Body),
				))
				msg.ParseMode = tgbotapi.ModeHTML
				for _, userTrackedProfile := range usersTrackedProfile {
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
	}
}

func (n *ProfileNotificator) AddTrackingProfile(userID uint, profileID uint) error {
	profile, err := n.GetUserProfile(profileID)
	if err != nil {
		return err
	}

	lastComment, err := n.Shiki.GetComments(profile.ID, shikitypes.TypeUser, 1, 1, true)
	if err != nil {
		return err
	}
	lastCommentID := uint(0)
	if len(lastComment) != 0 {
		lastCommentID = lastComment[0].ID
	}

	n.Database.Create(&models.TrackedProfile{
		UserID:        userID,
		ProfileID:     profileID,
		LastCommentID: lastCommentID,
	})

	return nil
}

func (n *ProfileNotificator) DeleteTrackingProfile(trackingProfileID uint, userID uint) error {
	return n.Database.Where("profile_id = ? AND user_id = ?", trackingProfileID, userID).Delete(&models.TrackedProfile{}).Error
}

func (n *ProfileNotificator) IsUserTrackingProfile(userID uint, profileID uint) bool {
	var trackedProfile models.TrackedProfile
	if err := n.Database.Where("user_id = ? AND profile_id = ?", userID, profileID).First(&trackedProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		panic(err)
	}
	return true
}

func (n *ProfileNotificator) GetUserProfile(id uint) (*shikitypes.UserProfile, error) {
	n.Mu.Lock()
	userProfile, ok := n.CachedProfiles[id]
	n.Mu.Unlock()
	if ok {
		return userProfile, nil
	}

	userProfile, err := n.Shiki.GetUserProfile(id)
	if err != nil {
		return nil, err
	}

	n.Mu.Lock()
	n.CachedProfiles[id] = userProfile
	n.Mu.Unlock()

	return userProfile, nil
}

func (n *ProfileNotificator) GetUserProfileByNickname(nickname string) (*shikitypes.UserProfile, error) {
	n.Mu.Lock()
	userProfile, ok := n.CachedProfilesByNickname[nickname]
	n.Mu.Unlock()
	if ok {
		return userProfile, nil
	}

	userProfile, err := n.Shiki.GetUserProfileByNickname(nickname)
	if err != nil {
		return nil, err
	}

	n.Mu.Lock()
	n.CachedProfilesByNickname[userProfile.Nickname] = userProfile
	n.Mu.Unlock()

	return userProfile, nil
}
