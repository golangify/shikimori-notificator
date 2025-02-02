package bbcode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var bbImageRegexp = regexp.MustCompile(`\[image=(\d+)\]`)

func (p *BBCodeParser) parseBBImage(text string) string {
	for {
		matches := bbImageRegexp.FindAllStringSubmatch(text, 1)
		if len(matches) == 0 {
			break
		}
		for _, match := range matches {
			imageID, _ := strconv.ParseUint(match[1], 10, 32)
			if _, err := p.shikiDB.GetImage(uint(imageID)); err == nil {
				text = strings.ReplaceAll(text, fmt.Sprint("[image=", imageID, "]"), fmt.Sprint("/image_", imageID))
			} else {
				text = strings.ReplaceAll(text, fmt.Sprint("[image=", imageID, "]"), fmt.Sprint("<del>[image=", imageID, "]</del>"))
			}
		}
	}
	return text
}
