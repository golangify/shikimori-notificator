package bbcode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var bbRepliesRegexp = regexp.MustCompile(`\[replies=((?:\d+|,|\s)*)\]`)

func (p *BBCodeParser) parseReplies(text string) string {
	fmt.Println(bbRepliesRegexp.FindAllStringSubmatch(text, -1))

	for _, match := range bbRepliesRegexp.FindAllStringSubmatch(text, -1) {
		var replies string
		for _, replyCommentIDstr := range strings.Split(match[1], ",") {
			replyCommentIDstr = strings.TrimSpace(replyCommentIDstr)
			replyCommentID, _ := strconv.ParseUint(replyCommentIDstr, 10, 32)
			if replyCommentID == 0 {
				continue
			}
			comment, err := p.shikiDB.GetComment(uint(replyCommentID))
			if err != nil {
				continue
			}
			replies += fmt.Sprint("[comment=", comment.ID, "]")
		}
		text = strings.ReplaceAll(text, match[0], replies)
	}
	return text
}
