package profilenotificator

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"shikimori-notificator/models"
	commentconstructor "shikimori-notificator/view/constructors/comment"
	"shikimori-notificator/workers/filter"
	shikidb "shikimori-notificator/workers/shiki-db"
	"slices"
	"strconv"
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
	ShikiDB  *shikidb.ShikiDB

	commentConstructor *commentconstructor.CommentConstructor
	filter             *filter.Filter

	ticker *time.Ticker
}

func NewProfileNotificator(bot *tgbotapi.BotAPI, shiki *shikimori.Client, shikidb *shikidb.ShikiDB, database *gorm.DB, filter *filter.Filter, commentConstructor *commentconstructor.CommentConstructor) *ProfileNotificator {
	n := &ProfileNotificator{
		Shiki:    shiki,
		Bot:      bot,
		Database: database,
		ShikiDB:  shikidb,

		commentConstructor: commentConstructor,

		filter: filter,
	}

	return n
}

func (n *ProfileNotificator) Run() {
	n.ticker = time.NewTicker(time.Minute)
	for range n.ticker.C {
		err := n.notifyProfiles()
		if err != nil {
			log.Println(err)
		}
		err = n.notifyPosting()
		if err != nil {
			log.Println(err)
		}
	}
}

func (n *ProfileNotificator) GetLast20PostedCommentsID(username string) ([]uint, error) {
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
	profile, err := n.ShikiDB.GetProfile(profileID)
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

	lastPostedComment, err := n.GetLast20PostedCommentsID(profile.Nickname)
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
