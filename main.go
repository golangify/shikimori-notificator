package main

import (
	"log"
	"shikimori-notificator/config"
	updatehandler "shikimori-notificator/handlers/update"
	"shikimori-notificator/models"
	"shikimori-notificator/view/constructors/bbcode"
	commentconstructor "shikimori-notificator/view/constructors/comment"
	topicconstructor "shikimori-notificator/view/constructors/topic"
	"shikimori-notificator/workers/filter"
	profilenotificator "shikimori-notificator/workers/profile-notificator"
	shikidb "shikimori-notificator/workers/shiki-db"
	topicnotificator "shikimori-notificator/workers/topic-notificator"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	config, err := config.LoadFromJsonFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	db, err := gorm.Open(sqlite.Open(config.Database.DatabaseString), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	if err = db.AutoMigrate(&models.User{}, &models.TrackedTopic{}, &models.TrackedProfile{}); err != nil {
		log.Fatalln(err)
	}

	bot, err := tgbotapi.NewBotAPI(config.Telegram.BotApiToken)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("авторизован в боте %s", bot.Self.UserName)

	shikimori.UserAgent = config.Shikimori.UserAgent
	shiki, err := shikimori.NewClient(config.Shikimori.Cookie, config.Shikimori.XsrfToken)
	if err != nil {
		log.Fatalln(err)
	}
	if shiki.Me.ID != 0 {
		log.Printf("авторизован в профиле %s", shiki.Me.Nickname)
	}

	shikiDB := shikidb.NewShikiDB(shiki)
	filter := filter.NewFilter()
	bbCodeParser := bbcode.NewBBCodeParser(shikiDB)

	commentsConstructor := commentconstructor.NewCommentConstructor(bbCodeParser)
	topicConstructor := topicconstructor.NewTopicConstructor(bbCodeParser)

	topicNotificator := topicnotificator.NewTopicNotificator(bot, shiki, db, shikiDB, filter, commentsConstructor)
	go topicNotificator.Run()

	profileNotificator := profilenotificator.NewProfileNotificator(bot, shiki, shikiDB, db, filter, commentsConstructor)
	go profileNotificator.Run()

	uh := updatehandler.New(bot, shiki, db, shikiDB, topicNotificator, profileNotificator, topicConstructor)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		go func(u tgbotapi.Update) {
			uh.Process(&u)
		}(update)
	}
}
