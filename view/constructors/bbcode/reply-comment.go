package bbcode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	shikimori "github.com/golangify/go-shiki-api"
)

var bbReplyCommentRegexp = regexp.MustCompile(`\[comment=(\d+);(\d+)\]`)

func (p *BBCodeParser) parseReplyComment(text string) string {
	for _, match := range bbReplyCommentRegexp.FindAllStringSubmatch(text, 1) {
		replyCommentID, _ := strconv.ParseUint(match[1], 10, 32)
		// userID, _ := strconv.ParseUint(match[2], 10, 32)
		replyCommentBBCode := fmt.Sprint("[comment=", replyCommentID, ";", match[2], "]")
		var nickname string
		if comment, err := p.shikiDB.GetComment(uint(replyCommentID)); err == nil {
			nickname = comment.User.Nickname
		} else {
			nickname = err.Error()
		}
		text = strings.ReplaceAll(
			text,
			replyCommentBBCode,
			fmt.Sprint(
				"<a href='", shikimori.ShikiSchema, "://", shikimori.ShikiDomain, "/comments/", replyCommentID, "'>@", nickname, "</a>",
			),
		)
	}
	return text
}
