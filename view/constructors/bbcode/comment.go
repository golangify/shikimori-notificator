package bbcode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	shikimori "github.com/golangify/go-shiki-api"
)

var bbCommentRegexp = regexp.MustCompile(`\[comment=(\d+);(\d+)\]`)

func (p *BBCodeParser) parseBBComment(text string) string {
	for {
		matches := bbCommentRegexp.FindAllStringSubmatch(text, 1)
		if len(matches) == 0 {
			break
		}
		for _, match := range matches {
			commentID, _ := strconv.ParseUint(match[1], 10, 32)
			userID, _ := strconv.ParseUint(match[2], 10, 32)
			var nickname string
			if user, err := p.shikiDB.GetProfile(uint(userID)); err == nil {
				nickname = user.Nickname
			} else {
				nickname = err.Error()
			}
			text = strings.ReplaceAll(
				text,
				fmt.Sprint("[comment=", commentID, ";", userID, "]"),
				fmt.Sprint(
					"<a href='", shikimori.ShikiSchema, "://", shikimori.ShikiDomain, "/", nickname, "'>", nickname, "</a> ",
					"<a href='", shikimori.ShikiSchema, "://", shikimori.ShikiDomain, "/comments/", commentID, "'>#</a>",
				),
			)
		}
	}
	return text
}
