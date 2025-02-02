package bbcode

import (
	shikidb "shikimori-notificator/workers/shiki-db"
	"strings"
)

type BBCodeParser struct {
	shikiDB *shikidb.ShikiDB
}

func NewBBCodeParser(shikidb *shikidb.ShikiDB) *BBCodeParser {
	return &BBCodeParser{
		shikiDB: shikidb,
	}
}

func (p *BBCodeParser) Parse(text string) string {
	text = strings.TrimSpace(text)
	runeText := []rune(text)
	if len(runeText) > 3900 {
		runeText = runeText[:3900]
	}
	text = string(runeText)
	text = p.parseSingleTags(text)
	text = p.parseDoubleTags(text)
	return text
}

func (p *BBCodeParser) parseSingleTags(text string) string {
	text = p.parseReplyComment(text)
	text = p.parseImage(text)
	text = p.parseUser(text)
	text = p.parseTopic(text)
	return text
}

func (p *BBCodeParser) parseDoubleTags(text string) string {
	text = p.parseBold(text)
	text = p.parseItalic(text)
	text = p.parseUnderline(text)
	text = p.parseSpoiler(text)
	return text
}
