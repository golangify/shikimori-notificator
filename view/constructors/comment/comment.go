package commentconstructor

import (
	"fmt"
	"html"
	"shikimori-notificator/view/constructors/bbcode"
	"strconv"

	shikimori "github.com/golangify/go-shiki-api"
	shikitypes "github.com/golangify/go-shiki-api/types"
)

type CommentConstructor struct {
	bbCodeParser *bbcode.BBCodeParser
}

func NewCommentConstructor(bbCodeParser *bbcode.BBCodeParser) *CommentConstructor {
	return &CommentConstructor{
		bbCodeParser: bbCodeParser,
	}
}

func (c *CommentConstructor) Profile(comment *shikitypes.Comment, profile *shikitypes.UserProfile) string {
	result := fmt.Sprint(
		"[user=", comment.User.ID, "] в профиле [user=", comment.CommentableID, "]\n\n",
		html.EscapeString(comment.Body),
	)
	result = c.bbCodeParser.Parse(result)
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
	result = c.bbCodeParser.Parse(result)
	return result
}
