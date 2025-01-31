package topicnotificator

import (
	"log"
	"shikimori-notificator/models"
	commentconstructor "shikimori-notificator/view/constructors/comment"
	"shikimori-notificator/workers/filter"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	shikitypes "github.com/golangify/go-shiki-api/types"
	"gorm.io/gorm"
)

type TopicNotificator struct {
	Shiki    *shikimori.Client
	Bot      *tgbotapi.BotAPI
	Database *gorm.DB

	CachedTopics map[uint]*shikitypes.Topic
	Mu           sync.Mutex

	Filter *filter.Filter
}

func NewTopicNotificator(shiki *shikimori.Client, bot *tgbotapi.BotAPI, database *gorm.DB, filter *filter.Filter) *TopicNotificator {
	ntfctr := &TopicNotificator{
		Shiki:    shiki,
		Bot:      bot,
		Database: database,

		Filter: filter,
	}

	ntfctr.CachedTopics = make(map[uint]*shikitypes.Topic)

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
			topic, err := n.GetTopic(trackedTopic.TopicID)
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
				msg := tgbotapi.NewMessage(0, commentconstructor.TopicToMessageText(&newComment, topic))
				msg.ParseMode = tgbotapi.ModeHTML
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
	topic, err := n.GetTopic(topicID)
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

func (n *TopicNotificator) GetTopic(id uint) (*shikitypes.Topic, error) {
	n.Mu.Lock()
	topic, ok := n.CachedTopics[id]
	n.Mu.Unlock()
	if ok {
		return topic, nil
	}

	topic, err := n.Shiki.GetTopic(id)
	if err != nil {
		return nil, err
	}

	n.Mu.Lock()
	n.CachedTopics[id] = topic
	n.Mu.Unlock()

	return topic, nil
}
