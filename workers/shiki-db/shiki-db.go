package shikidb

import (
	"errors"
	"fmt"
	"regexp"
	"shikimori-notificator/workers/shiki-db/cacher"

	shikimori "github.com/golangify/go-shiki-api"
	shikitypes "github.com/golangify/go-shiki-api/types"
)

type ShikiDB struct {
	shiki  *shikimori.Client
	cacher *cacher.Cacher
}

func NewShikiDB(shiki *shikimori.Client) *ShikiDB {
	return &ShikiDB{
		shiki:  shiki,
		cacher: cacher.NewCacher(),
	}
}

func (d *ShikiDB) NumCached() uint {
	return d.cacher.NumCached()
}

func (d *ShikiDB) GetComment(commentID uint) (*shikitypes.Comment, error) {
	comment := d.cacher.GetComment(commentID)
	if comment != nil {
		return comment, nil
	}

	comment, err := d.shiki.GetComment(commentID)
	if err != nil {
		return nil, err
	}

	d.cacher.SetComment(comment.ID, *comment)

	return comment, nil
}

func (d *ShikiDB) GetTopic(topicID uint) (*shikitypes.Topic, error) {
	topic := d.cacher.GetTopic(topicID)
	if topic != nil {
		return topic, nil
	}

	topic, err := d.shiki.GetTopic(topicID)
	if err != nil {
		return nil, err
	}

	d.cacher.SetTopic(topic.ID, *topic)

	return topic, nil
}

func (d *ShikiDB) GetProfile(profileID uint) (*shikitypes.UserProfile, error) {
	profile := d.cacher.GetProfile(profileID)
	if profile != nil {
		return profile, nil
	}

	profile, err := d.shiki.GetUserProfile(profileID)
	if err != nil {
		return nil, err
	}

	d.cacher.SetProfile(profileID, *profile)

	return profile, nil
}

func (d *ShikiDB) GetProfileByNickname(nickname string) (*shikitypes.UserProfile, error) {
	profile := d.cacher.GetProfileByNickname(nickname)
	if profile != nil {
		return profile, nil
	}

	profile, err := d.shiki.GetUserProfileByNickname(nickname)
	if err != nil {
		return nil, err
	}

	d.cacher.SetProfileByNickname(nickname, *profile)

	return profile, nil
}

func (d *ShikiDB) ClearCache() uint {
	return d.cacher.Clear()
}

var imageLinkRegexp = regexp.MustCompile(`((?:http|https):\/\/[a-z]+\.shikimori\.one\/system\/user_images_h\/original\/[a-z0-9]+\/[a-z0-9]+\.jpg)`)

func (d *ShikiDB) GetImageLink(imageID uint) (*string, error) {
	imageLink := d.cacher.GetImage(imageID)
	if imageLink != nil {
		return imageLink, nil
	}

	data, err := d.shiki.PreviewComment(fmt.Sprint("[image=", imageID, "]"))
	if err != nil {
		return nil, err
	}

	matches := imageLinkRegexp.FindAllString(string(data), -1)
	if len(matches) == 0 {
		return nil, errors.New("изображение не найдено")
	}

	imageLink = &matches[0]
	d.cacher.SetImage(imageID, *imageLink)

	return imageLink, nil
}
