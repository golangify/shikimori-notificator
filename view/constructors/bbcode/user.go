package bbcode

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	shikimori "github.com/golangify/go-shiki-api"
)

var bbUserRegexp = regexp.MustCompile(`\[user=(\d+)\]`)

func (p *BBCodeParser) parseUser(text string) string {
	for _, match := range bbUserRegexp.FindAllStringSubmatch(text, -1) {
		userID, _ := strconv.ParseUint(match[1], 10, 32)
		userBBCode := fmt.Sprint("[user=", userID, "]")
		user, err := p.shikiDB.GetProfile(uint(userID))
		var nickname string
		if err == nil {
			nickname = user.Nickname
		} else {
			nickname = err.Error()
		}
		text = strings.ReplaceAll(text, userBBCode, fmt.Sprint("<a href='", shikimori.ShikiSchema, "://", shikimori.ShikiDomain, "/", url.QueryEscape(nickname), "'>", html.EscapeString(nickname), "</a>"))
	}

	return text
}
