package bbcode

import (
	"fmt"
	"regexp"
	"strings"
)

func (p *BBCodeParser) parseSpoiler(text string) string {
	text = p.parseInlineSpoiler(text)
	text = p.parseSpoilerBlock(text)
	text = p.parseSpoilerWithTitle(text)
	return text
}

var bbInlineSpoiler = regexp.MustCompile(`\|\|(.+)\|\|`) // содержимое строго в 1 строку

func (p *BBCodeParser) parseInlineSpoiler(text string) string {
	for _, match := range bbInlineSpoiler.FindAllStringSubmatch(text, -1) {
		inlineSpoilerBody := match[1]
		text = strings.ReplaceAll(text, match[0], fmt.Sprint("<span class='tg-spoiler'>", inlineSpoilerBody, "</span>"))
	}
	return text
}

var bbSpoilerBlock = regexp.MustCompile(`\[spoiler_block\]((?s).+)\[/spoiler_block]`)

func (p *BBCodeParser) parseSpoilerBlock(text string) string {
	for _, match := range bbSpoilerBlock.FindAllStringSubmatch(text, -1) {
		spoilerBlockBody := match[1]
		text = strings.ReplaceAll(text, match[0], fmt.Sprint("<span class='tg-spoiler'>", spoilerBlockBody, "</span>"))
	}
	return text
}

var bbSpoilerWithTitle = regexp.MustCompile(`\[spoiler=(.+)\]((?s).+)\[/spoiler\]`)

func (p *BBCodeParser) parseSpoilerWithTitle(text string) string {
	for _, match := range bbSpoilerWithTitle.FindAllStringSubmatch(text, -1) {
		spoilerTitle := match[1]
		spoilerBlockBody := match[2]
		text = strings.ReplaceAll(text, match[0], fmt.Sprint("(", spoilerTitle, ")<span class='tg-spoiler'>", spoilerBlockBody, "</span>"))
	}
	return text
}
