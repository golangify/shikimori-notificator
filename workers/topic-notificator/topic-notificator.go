package topicnotificator

import (
	"fmt"
	"html"
	"log"
	"shikimori-notificator/models"
	"strconv"
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

	ticker *time.Ticker
}

func New(shiki *shikimori.Client, bot *tgbotapi.BotAPI, database *gorm.DB) *TopicNotificator {
	ntfctr := &TopicNotificator{
		Shiki:    shiki,
		Bot:      bot,
		Database: database,
	}

	ntfctr.CachedTopics = make(map[uint]*shikitypes.Topic)

	return ntfctr
}

func (n *TopicNotificator) Run() {
	n.ticker = time.NewTicker(time.Minute)
	for range n.ticker.C {
		var trackedTopics []models.TrackedTopic
		n.Database.Find(&trackedTopics).Distinct("topic_id")
		t := time.NewTicker(time.Second * 2)
		for _, topic := range trackedTopics {
			<-t.C
			lastComments, err := n.Shiki.GetComments(topic.TopicID, shikitypes.TypeTopic, 1, 20, true)
			if err != nil {
				log.Println(err)
				continue
			}
			var newComments []shikitypes.Comment
			for _, comment := range lastComments {
				if comment.ID > topic.LastCommentID {
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
			err = n.Database.Preload("User").Find(&usersTrackedTopic, topic.ID).Error
			if err != nil {
				log.Println(err)
				continue
			}
			for _, newComment := range newComments {
				topic, err := n.GetTopic(newComment.CommentableID)
				if err != nil {
					log.Println(err)
					continue
				}
				msg := tgbotapi.NewMessage(0, fmt.Sprintf(
					"<a href='%s'>%s</a> в топике <a href='%s'>%s</a>\n\n%s",
					newComment.User.URL, html.EscapeString(newComment.User.Nickname),
					shikimori.ShikiSchema+"://"+shikimori.ShikiDomain+topic.Forum.URL+"/"+strconv.FormatUint(uint64(topic.ID), 10), topic.TopicTitle,
					html.EscapeString(newComment.Body),
				))
				msg.ParseMode = tgbotapi.ModeHTML
				for _, userTrackedTopic := range usersTrackedTopic {
					<-t.C
					msg.BaseChat.ChatID = userTrackedTopic.User.TgID
					_, err = n.Bot.Send(msg)
					if err != nil {
						log.Println(err)
					}
				}
			}
			err = n.Database.Model(&models.TrackedTopic{}).Where("topic_id = ?", topic.TopicID).UpdateColumn("last_comment_id", newComments[len(newComments)-1].ID).Error
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// привязать новый отслеживаемый топик к пользователю
func (n *TopicNotificator) AddTrackingTopic(userID uint, topicID uint) error {
	// проверяем существует ли топик
	topic, err := n.GetTopic(topicID)
	if err != nil {
		return err
	}
	// проверяем существует ли пользователь
	var user models.User
	if err = n.Database.First(&user, userID).Error; err != nil {
		return err
	}

	// получаем последний комментарий топика
	lastComment, err := n.Shiki.GetComments(topicID, shikitypes.TypeTopic, 1, 1, true)
	if err != nil {
		panic(err)
	}
	lastCommentID := uint(0)
	if len(lastComment) != 0 {
		lastCommentID = lastComment[0].ID
	}
	// вносим в базу данных
	n.Database.Create(&models.TrackedTopic{
		UserID:        user.ID,
		TopicID:       topic.ID,
		LastCommentID: lastCommentID,
	})
	return nil
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
