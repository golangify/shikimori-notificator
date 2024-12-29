package commandhandler

import (
	"regexp"
	"shikimori-notificator/models"
	topicnotificator "shikimori-notificator/workers/topic-notificator"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	"gorm.io/gorm"
)

type command struct {
	Level            uint             // уровень прав необходимый пользователю для доступа
	Name             string           // topic
	Usage            string           // /topic [id]
	ActivatorRegexps []*regexp.Regexp // []*regexp.Regexp{regexp.MustCompile("^\/topic$")}
	Regexp           *regexp.Regexp   // regexp.MustCompile("^\/topic (\d+)")
	Description      string           // получить топик по id
	Function         commandFunction
}

type commandFunction func(update *tgbotapi.Update, user *models.User, args []string)

func (c *command) Help() string {
	helpText := c.Usage + " - " + c.Description
	if c.Function == nil {
		helpText += " (команда временно недоступна)"
	}
	return helpText
}

type CommandHandler struct {
	Bot              *tgbotapi.BotAPI
	Shiki            *shikimori.Client
	TopicNotificator *topicnotificator.TopicNotificator

	Database *gorm.DB

	commands []*command
}

func NewCommandHandler(bot *tgbotapi.BotAPI, shiki *shikimori.Client, topicNotificator *topicnotificator.TopicNotificator, db *gorm.DB) *CommandHandler {
	h := &CommandHandler{
		Bot:              bot,
		Shiki:            shiki,
		TopicNotificator: topicNotificator,
		Database:         db,
	}

	h.commands = []*command{
		{
			Name:        "help",
			Usage:       "/help",
			Regexp:      regexp.MustCompile(`^\/help$`),
			Description: "помощь",
			Function:    h.Help,
		},
		{
			Name:  "topic",
			Usage: "/topic [id]",
			ActivatorRegexps: []*regexp.Regexp{
				regexp.MustCompile(`^\/topic$`),
			},
			Regexp:      regexp.MustCompile(`^\/topic (\d+)$`),
			Description: "получить топик по id",
			Function:    h.Topic,
		},
		{
			Name:        "topics",
			Usage:       "/topics",
			Regexp:      regexp.MustCompile(`^\/topics$`),
			Description: "получить список моих отслеживаемых топиков",
			Function:    h.Topics,
		},
		{
			Level:       3,
			Name:        "debug",
			Usage:       "/debug",
			Regexp:      regexp.MustCompile(`^\/debug$`),
			Description: "активировать debug режим",
			Function:    h.Debug,
		},
		{
			Level: 3,
			Name:  "disablecommand",
			Usage: "/disablecommand [command]",
			ActivatorRegexps: []*regexp.Regexp{
				regexp.MustCompile(`^\/disablecommand$`),
			},
			Regexp:      regexp.MustCompile(`^\/disablecommand ([a-z]+)$`),
			Description: "отключить команду в боте",
			Function:    h.Disablecommand,
		},
		{
			Level: 3,
			Name:  "enablecommand",
			Usage: "/enablecommand [command]",
			ActivatorRegexps: []*regexp.Regexp{
				regexp.MustCompile(`^\/enablecommand$`),
			},
			Regexp:      regexp.MustCompile(`^\/enablecommand ([a-z]+)$`),
			Description: "включить отключенную команду в боте",
			Function:    h.Enablecommand,
		},
		{
			Level:       3,
			Name:        "clearcache",
			Usage:       "/clearcache",
			Regexp:      regexp.MustCompile(`^\/clearcache$`),
			Description: "очистить кэш объектов",
			Function:    h.Clearcache,
		},
	}

	return h

}

func (h *CommandHandler) Process(update *tgbotapi.Update, user *models.User) {
	for _, cmd := range h.commands {
		if cmd.Level > user.Level {
			continue // пропускаем команду, т.к. у пользователя недостаточно прав на её использование
		}
		// сначала ищем полное правильное совпадение регекспы функции
		if cmd.Regexp.MatchString(update.Message.Text) {
			if cmd.Function == nil {
				h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда недоступна."))
				return
			}
			cmd.Function(update, user, cmd.Regexp.FindAllStringSubmatch(update.Message.Text, -1)[0])
			return
		}
		// если не нашли полное совпадение, то ищем по минимально-необходимому для определения того, чего хотел пользователь
		for _, activatorRegexp := range cmd.ActivatorRegexps {
			if activatorRegexp.MatchString(update.Message.Text) {
				h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, cmd.Help()))
				return
			}
		}
	}
	h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда."))
}
