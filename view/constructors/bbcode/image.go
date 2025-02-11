package bbcode

import (
	"fmt"
	"regexp"
	"strings"
)

var bbImageRegexp = regexp.MustCompile(`\[(?:image|poster)=(\d+).+?[^\[]]`)

func (p *BBCodeParser) parseImage(text string) string {
	for _, match := range bbImageRegexp.FindAllStringSubmatch(text, -1) {
		text = strings.ReplaceAll(text, match[0], fmt.Sprint("/image", match[1]))
	}
	return text
}
