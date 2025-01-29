package commentconstructor

import (
	"fmt"
	"html"
	"strconv"

	shikimori "github.com/golangify/go-shiki-api"
	shikitypes "github.com/golangify/go-shiki-api/types"
)

func ProfileToMessageText(c *shikitypes.Comment, commentableUser *shikitypes.UserProfile) string {
	result := fmt.Sprintf(
		"<a href='%s'>%s</a> в профиле <a href='%s'>%s</a>\n\n%s",
		c.User.URL, html.EscapeString(c.User.Nickname),
		commentableUser.URL, html.EscapeString(commentableUser.Nickname),
		html.EscapeString(c.Body),
	)
	runeResult := []rune(result)
	if len(runeResult) > 4096 {
		runeResult = runeResult[:4096]
	}
	result = string(runeResult)
	return result
}

func TopicToMessageText(c *shikitypes.Comment, commentableTopic *shikitypes.Topic) string {
	result := fmt.Sprintf(
		"<a href='%s'>%s</a> в топике <a href='%s'>%s</a>\n\n%s",
		c.User.URL, html.EscapeString(c.User.Nickname),
		shikimori.ShikiSchema+"://"+shikimori.ShikiDomain+commentableTopic.Forum.URL+"/"+strconv.FormatUint(uint64(commentableTopic.ID), 10), commentableTopic.TopicTitle,
		html.EscapeString(c.Body),
	)
	runeResult := []rune(result)
	if len(runeResult) > 4096 {
		runeResult = runeResult[:4096]
	}
	result = string(runeResult)
	return result
}
