package shikidb

import (
	"shikimori-notificator/workers/cacher"

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
