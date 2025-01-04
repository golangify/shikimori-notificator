package topicconstructor

import (
	"fmt"
	"html"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	shikitypes "github.com/golangify/go-shiki-api/types"
)

func ToMessageText(t *shikitypes.Topic) string {
	messageText := fmt.Sprintf("<a href='%s'>%s</a>\n\n%s",
		shikimori.ShikiSchema+"://"+shikimori.ShikiDomain+t.Forum.URL+"/"+fmt.Sprint(t.ID),
		t.TopicTitle,
		html.EscapeString(t.Body),
	)
	runeMessageText := []rune(messageText)
	if len(runeMessageText) > 3900 {
		runeMessageText = runeMessageText[:3900]
	}
	return string(runeMessageText)
}

func ToInlineKeyboard(t *shikitypes.Topic, isTopicTracking bool) *tgbotapi.InlineKeyboardMarkup {
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
