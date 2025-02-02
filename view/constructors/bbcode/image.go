package bbcode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var bbDeletedImageRegexp = regexp.MustCompile(`<del>\[image=(\d+)\]</del>`)
var bbImageRegexp = regexp.MustCompile(`\[image=(\d+)\]`)

func (p *BBCodeParser) parseImage(text string) string {
	for _, match := range bbDeletedImageRegexp.FindAllStringSubmatch(text, -1) {
		imageID, _ := strconv.ParseUint(match[1], 10, 32)
		imageBBCode := fmt.Sprint("[deleted_image=", imageID, "]")
		text = strings.ReplaceAll(text, imageBBCode, fmt.Sprint("<del>", imageBBCode, "</del>"))
	}
	for _, match := range bbImageRegexp.FindAllStringSubmatch(text, -1) {
		imageID, _ := strconv.ParseUint(match[1], 10, 32)
		imageBBCode := fmt.Sprint("[image=", imageID, "]")

		if _, err := p.shikiDB.GetImageLink(uint(imageID)); err == nil {
			text = strings.ReplaceAll(text, imageBBCode, fmt.Sprint("[/image", imageID, "]"))
		} else {
			text = strings.ReplaceAll(text, imageBBCode, fmt.Sprint("<del>", imageBBCode, "</del>"))
		}
	}
	text = strings.ReplaceAll(text, "[deleted_image=", "[image=")
	return text
}
