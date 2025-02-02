package bbcode

import (
	"fmt"
	"regexp"
	"strings"
)

var bbBoldRegexp = regexp.MustCompile(`\[b\]((?s).+)\[/b\]`)

func (p *BBCodeParser) parseBold(text string) string {
	for _, match := range bbBoldRegexp.FindAllStringSubmatch(text, -1) {
		boldTextBody := match[1]
		text = strings.ReplaceAll(text, match[0], fmt.Sprint("<b>", boldTextBody, "</b>"))
	}
	return text
}
