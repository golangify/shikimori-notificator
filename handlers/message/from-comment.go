package messagehandler

import (
	"fmt"
	"shikimori-notificator/models"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
	shikitypes "github.com/golangify/go-shiki-api/types"
)

func (h *MessageHandler) FromComment(update *tgbotapi.Update, user *models.User, args []string) {
	commentID, _ := strconv.ParseUint(args[1], 10, 32)
	comment, err := h.Shiki.GetComment(uint(commentID))
	if err != nil {
		if err == shikimori.ErrNotFound {
			h.Bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Комментарий с id ", commentID, " не найден.")))
			return
		}
	}
	if comment.CommentableType == shikitypes.TypeUser {
		msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("<code>/profile ", comment.CommentableID, "</code>"))
		msg.ParseMode = tgbotapi.ModeHTML
		h.Bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("<code>/topic ", comment.CommentableID, "</code>"))
		msg.ParseMode = tgbotapi.ModeHTML
		h.Bot.Send(msg)
	}
}
