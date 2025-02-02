package bbcode

import (
	"fmt"
	"regexp"
	"strings"
)

var bbItalicRegexp = regexp.MustCompile(`\[i\]((?s).+)\[/i\]`)

func (p *BBCodeParser) parseItalic(text string) string {
	for _, match := range bbItalicRegexp.FindAllStringSubmatch(text, -1) {
		italicTextBody := match[1]
		text = strings.ReplaceAll(text, match[0], fmt.Sprint("<i>", italicTextBody, "</i>"))
	}
	return text
}
