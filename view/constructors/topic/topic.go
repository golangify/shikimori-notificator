package topicconstructor

import (
	"fmt"
	"html"
	"shikimori-notificator/view/constructors/bbcode"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	shikitypes "github.com/golangify/go-shiki-api/types"
)

type TopicConstructor struct {
	bbCodeParser *bbcode.BBCodeParser
}

func NewTopicConstructor(bbCodeParser *bbcode.BBCodeParser) *TopicConstructor {
	return &TopicConstructor{
		bbCodeParser: bbCodeParser,
	}
}

func (p *TopicConstructor) Text(t *shikitypes.Topic) string {
	messageText := fmt.Sprintf("<a href='%s'>%s</a>\n\n%s",
		shikimori.ShikiSchema+"://"+shikimori.ShikiDomain+t.Forum.URL+"/"+fmt.Sprint(t.ID),
		t.TopicTitle,
		html.EscapeString(t.Body),
	)
	messageText = p.bbCodeParser.Parse(messageText)
	return messageText
}

func (p *TopicConstructor) InlineKeyboard(t *shikitypes.Topic, isTopicTracking bool) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	if isTopicTracking {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("прекратить отслеживать топик", fmt.Sprint("delete_topic_from_tracking ", t.ID)),
		))
	} else {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("отслеживать топик", fmt.Sprint("add_topic_to_tracking ", t.ID)),
		))
	}
	return &keyboard
}
