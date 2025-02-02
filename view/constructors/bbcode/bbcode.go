package bbcode

import (
	shikidb "shikimori-notificator/workers/shiki-db"
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
	runeText := []rune(text)
	if len(runeText) > 3900 {
		runeText = runeText[:3900]
	}
	text = string(runeText)
	text = p.parseSingleTags(text)
	return text
}

func (p *BBCodeParser) parseSingleTags(text string) string {
	text = p.parseBBComment(text)
	text = p.parseBBImage(text)
	return text
}
