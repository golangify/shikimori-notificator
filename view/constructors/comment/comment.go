package commentconstructor

import (
	"fmt"
	"html"
	"shikimori-notificator/view/constructors/bbcode"
	shikidb "shikimori-notificator/workers/shiki-db"
	"strconv"

	shikimori "github.com/golangify/go-shiki-api"
	shikitypes "github.com/golangify/go-shiki-api/types"
)

type CommentConstructor struct {
	bbCodeParser *bbcode.BBCodeParser
}

func NewCommentConstructor(shikidb *shikidb.ShikiDB) *CommentConstructor {
	return &CommentConstructor{
		bbCodeParser: bbcode.NewBBCodeParser(shikidb),
	}
}

func (c *CommentConstructor) Profile(comment *shikitypes.Comment, profile *shikitypes.UserProfile) string {
	result := fmt.Sprintf(
		"<a href='%s'>%s</a> в профиле <a href='%s'>%s</a> <a href='%s://%s/comments/%d'>#</a>\n\n%s",
		comment.User.URL, html.EscapeString(comment.User.Nickname),
		profile.URL, html.EscapeString(profile.Nickname),
		shikimori.ShikiSchema, shikimori.ShikiDomain, comment.ID,
		html.EscapeString(comment.Body),
	)
	runeResult := []rune(result)
	if len(runeResult) > 4096 {
		runeResult = runeResult[:4096]
	}
	result = c.bbCodeParser.Parse(string(runeResult))
	return result
}

func (c *CommentConstructor) Topic(comment *shikitypes.Comment, commentableTopic *shikitypes.Topic) string {
	result := fmt.Sprintf(
		"<a href='%s'>%s</a> в топике <a href='%s'>%s</a> <a href='%s://%s/comments/%d'>#</a>\n\n%s",
		comment.User.URL, html.EscapeString(comment.User.Nickname),
		shikimori.ShikiSchema+"://"+shikimori.ShikiDomain+commentableTopic.Forum.URL+"/"+strconv.FormatUint(uint64(commentableTopic.ID), 10), commentableTopic.TopicTitle,
		shikimori.ShikiSchema, shikimori.ShikiDomain, comment.ID,
		html.EscapeString(comment.Body),
	)
	runeResult := []rune(result)
	if len(runeResult) > 4096 {
		runeResult = runeResult[:4096]
	}
	result = c.bbCodeParser.Parse(string(runeResult))
	return result
}
