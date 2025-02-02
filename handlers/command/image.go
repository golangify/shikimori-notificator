package commandhandler

import (
	"fmt"
	"net/http"
	"regexp"
	"shikimori-notificator/models"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	shikimori "github.com/golangify/go-shiki-api"
)

var removePrefixImageLinkRegexp = regexp.MustCompile(`(https:\/\/[a-z]+\.shikimori\.one\/)`)

func (h *CommandHandler) image(update *tgbotapi.Update, _ *models.User, args []string) {
	imageID, _ := strconv.ParseUint(args[1], 10, 32)
	imageLink, err := h.ShikiDB.GetImageLink(uint(imageID))
	if err != nil {
		panic(err)
	}

	*imageLink = removePrefixImageLinkRegexp.ReplaceAllString(*imageLink, "")
	resp, err := h.Shiki.MakeRequest(http.MethodGet, strings.TrimPrefix(*imageLink, shikimori.ShikiSchema+"://"+shikimori.ShikiDomain+"/"), nil, nil, nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	msg := tgbotapi.NewPhoto(update.FromChat().ID, tgbotapi.FileReader{
		Name:   fmt.Sprint("image-", imageID, ".jpg"),
		Reader: resp.Body,
	})
	msg.Caption = fmt.Sprint(
		"<code>[image=", imageID, "]</code>\n\n",
		"<code>[img]", *imageLink, "[/img]</code>",
	)
	msg.ParseMode = tgbotapi.ModeHTML

	h.Bot.Send(msg)
}
