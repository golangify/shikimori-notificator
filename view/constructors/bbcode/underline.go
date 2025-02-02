package bbcode

import (
	"fmt"
	"regexp"
	"strings"
)

var bbUnderlineRegexp = regexp.MustCompile(`\[u\]((?s).+)\[/u\]`)

func (p *BBCodeParser) parseUnderline(text string) string {
	for _, match := range bbUnderlineRegexp.FindAllStringSubmatch(text, -1) {
		underlineTextBody := match[1]
		text = strings.ReplaceAll(text, match[0], fmt.Sprint("<u>", underlineTextBody, "</u>"))
	}
	return text
}
