package topicnotificator

import (
	"log"
	"shikimori-notificator/models"
	commentconstructor "shikimori-notificator/view/constructors/comment"
	"shikimori-notificator/workers/filter"
	shikidb "shikimori-notificator/workers/shiki-db"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	shikitypes "github.com/golangify/go-shiki-api/types"
	"gorm.io/gorm"
)

type TopicNotificator struct {
	Shiki              *shikimori.Client
	Bot                *tgbotapi.BotAPI
	Database           *gorm.DB
	ShikiDB            *shikidb.ShikiDB
	commentConstructor *commentconstructor.CommentConstructor

	Filter *filter.Filter
}

func NewTopicNotificator(bot *tgbotapi.BotAPI, shiki *shikimori.Client, database *gorm.DB, shikidb *shikidb.ShikiDB, filter *filter.Filter, commentConstructor *commentconstructor.CommentConstructor) *TopicNotificator {
	ntfctr := &TopicNotificator{
		Shiki:              shiki,
		Bot:                bot,
		Database:           database,
		ShikiDB:            shikidb,
		commentConstructor: commentConstructor,

		Filter: filter,
	}

	return ntfctr
}

func (n *TopicNotificator) Run() {
	t := time.NewTicker(time.Minute)
	for range t.C {
		var trackedTopics []models.TrackedTopic
		n.Database.Find(&trackedTopics).Order("last_comment_id").Distinct("topic_id")
		t := time.NewTicker(time.Second * 2)
		for _, trackedTopic := range trackedTopics {
			<-t.C
			topic, err := n.ShikiDB.GetTopic(trackedTopic.TopicID)
			if err != nil {
				log.Println(err)
				continue
			}
			lastComments, err := n.Shiki.GetComments(topic.ID, shikitypes.TypeTopic, 1, 20, true)
			if err != nil {
				log.Println(err)
				continue
			}
			var newComments []shikitypes.Comment
			for _, comment := range lastComments {
				if comment.ID > trackedTopic.LastCommentID {
					newComments = append(newComments, comment)
				}
			}
			if len(newComments) == 0 {
				continue
			}

			for i, j := 0, len(newComments)-1; i < j; i, j = i+1, j-1 {
				newComments[i], newComments[j] = newComments[j], newComments[i]
			}

			var usersTrackedTopic []models.TrackedTopic
			err = n.Database.Preload("User").Find(&usersTrackedTopic, "topic_id = ?", topic.ID).Error
			if err != nil {
				log.Println(err)
				continue
			}
			for _, newComment := range newComments {
				msg := tgbotapi.NewMessage(0, n.commentConstructor.Topic(&newComment, topic))
				msg.ParseMode = tgbotapi.ModeHTML
				msg.DisableWebPagePreview = true
				for _, userTrackedTopic := range usersTrackedTopic {
					if !n.Filter.Ok(newComment.ID, userTrackedTopic.User.ID) {
						continue
					}
					<-t.C
					msg.BaseChat.ChatID = userTrackedTopic.User.TgID
					_, err = n.Bot.Send(msg)
					if err != nil {
						log.Println(err)
					}
				}
			}
			err = n.Database.Model(&models.TrackedTopic{}).Where("topic_id = ?", topic.ID).UpdateColumn("last_comment_id", newComments[len(newComments)-1].ID).Error
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (n *TopicNotificator) AddTrackingTopic(userID uint, topicID uint) error {
	topic, err := n.ShikiDB.GetTopic(topicID)
	if err != nil {
		return err
	}

	lastComment, err := n.Shiki.GetComments(topicID, shikitypes.TypeTopic, 1, 1, true)
	if err != nil {
		panic(err)
	}
	lastCommentID := uint(0)
	if len(lastComment) != 0 {
		lastCommentID = lastComment[0].ID
	}
	n.Database.Create(&models.TrackedTopic{
		UserID:        userID,
		TopicID:       topic.ID,
		LastCommentID: lastCommentID,
	})
	return nil
}

func (n *TopicNotificator) DeleteTrackingTopic(userID uint, topicID uint) error {
	return n.Database.Where("topic_id = ? AND user_id = ?", topicID, userID).Delete(&models.TrackedTopic{}).Error
}

func (n *TopicNotificator) IsUserTrackingTopic(userID uint, topicID uint) bool {
	var trackedTopic models.TrackedTopic
	if err := n.Database.Where("user_id = ? AND topic_id = ?", userID, topicID).First(&trackedTopic).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		panic(err)
	}
	return true
}
