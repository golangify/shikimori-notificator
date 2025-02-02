package bbcode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	shikimori "github.com/golangify/go-shiki-api"
)

var bbTopicRegexp = regexp.MustCompile(`\[topic=(\d+)\]`)

func (p *BBCodeParser) parseTopic(text string) string {
	for _, match := range bbTopicRegexp.FindAllStringSubmatch(text, -1) {
		topicID, _ := strconv.ParseUint(match[1], 10, 32)
		topic, err := p.shikiDB.GetTopic(uint(topicID))
		var title string
		if err == nil {
			title = topic.TopicTitle
		} else {
			title = err.Error()
		}
		topicLink := fmt.Sprint(
			shikimori.ShikiSchema, "://", shikimori.ShikiDomain, topic.Forum.URL, "/", topic.ID,
		)
		text = strings.ReplaceAll(text, match[0], fmt.Sprint("<a href='", topicLink, "'>", title, "</a>"))
	}

	return text
}
