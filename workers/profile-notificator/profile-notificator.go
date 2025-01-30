package profilenotificator

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"shikimori-notificator/models"
	commentconstructor "shikimori-notificator/view/constructors/comment"
	topicnotificator "shikimori-notificator/workers/topic-notificator"
	"slices"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	shikitypes "github.com/golangify/go-shiki-api/types"
	"gorm.io/gorm"
)

var commentIDRregexp = regexp.MustCompile(`data-track_comment=\\"(\d+)\\"`)

type ProfileNotificator struct {
	Shiki    *shikimori.Client
	Bot      *tgbotapi.BotAPI
	Database *gorm.DB

	TopicNotificator *topicnotificator.TopicNotificator

	CachedProfiles           map[uint]*shikitypes.UserProfile
	CachedProfilesByNickname map[string]*shikitypes.UserProfile
	Mu                       sync.Mutex
	ticker                   *time.Ticker
}

func NewProfileNotificator(shiki *shikimori.Client, bot *tgbotapi.BotAPI, database *gorm.DB, topicNotificator *topicnotificator.TopicNotificator) *ProfileNotificator {
	n := &ProfileNotificator{
		Shiki:    shiki,
		Bot:      bot,
		Database: database,

		TopicNotificator: topicNotificator,

		CachedProfiles:           make(map[uint]*shikitypes.UserProfile),
		CachedProfilesByNickname: make(map[string]*shikitypes.UserProfile),
	}

	return n
}

func (n *ProfileNotificator) Run() {
	n.ticker = time.NewTicker(time.Second)
	filter := NewDuplicatorFilter()
	for range n.ticker.C {

		var trackedProfiles []models.TrackedProfile
		err := n.Database.Find(&trackedProfiles).Order("last_comment_id").Distinct("profile_id").Error
		if err != nil {
			log.Println(err)
			continue
		}
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
				for _, userTrackedProfile := range usersTrackedProfile {
					if filter.Check(newComment.ID, userTrackedProfile.UserID) {
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

		// проверка поcтинга
		trackedProfiles = []models.TrackedProfile{}
		if err = n.Database.Find(&trackedProfiles, "track_posting = ?", true).Order("last_posted_comment_id").Distinct("profile_id").Error; err != nil {
			log.Println(err)
			continue
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
				commment, err := n.Shiki.GetComment(newCommentID)
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
					if filter.Check(newComment.ID, userTrackedProfile.UserID) {
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
		t.Stop()
	}
}

func (n *ProfileNotificator) GetLast20PostedCommentIDs(username string) ([]uint, error) {
	var result []uint
	resp, err := n.Shiki.MakeRequest(http.MethodGet, "/"+username+"/comments", nil, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	htmlByteData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	matches := commentIDRregexp.FindAllStringSubmatch(string(htmlByteData), -1)
	for _, match := range matches {
		commentID, err := strconv.ParseUint(match[1], 10, 32)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(commentID))
	}
	slices.Reverse(result)
	return result, nil
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

	lastPostedComment, err := n.GetLast20PostedCommentIDs(profile.Nickname)
	if err != nil {
		return err
	}
	lastPostedCommentID := uint(0)
	if len(lastPostedComment) != 0 {
		lastPostedCommentID = lastPostedComment[len(lastPostedComment)-1]
	}

	n.Database.Create(&models.TrackedProfile{
		UserID:              userID,
		ProfileID:           profileID,
		LastCommentID:       lastCommentID,
		LastPostedCommentID: lastPostedCommentID,
	})

	return nil
}

func (n *ProfileNotificator) GetTrackingProfile(targetProfileID, userID uint) (*models.TrackedProfile, error) {
	var result models.TrackedProfile
	if err := n.Database.First(&result, "profile_id = ? AND user_id = ?", targetProfileID, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
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
