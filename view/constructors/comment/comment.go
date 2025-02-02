package commentconstructor

import (
	"fmt"
	"html"
	"shikimori-notificator/view/constructors/bbcode"

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
	result := fmt.Sprint(
		"[user=", comment.User.ID, "] в топике [topic=", comment.CommentableID, "]\n\n",
		html.EscapeString(comment.Body),
	)
	result = c.bbCodeParser.Parse(result)
	return result
}
