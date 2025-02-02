package cacher

import (
	cachestorage "shikimori-notificator/workers/shiki-db/cacher/storage"

	shikitypes "github.com/golangify/go-shiki-api/types"
)

type Cacher struct {
	cachedComments           *cachestorage.CacheStorage[uint, shikitypes.Comment]
	cachedProfiles           *cachestorage.CacheStorage[uint, shikitypes.UserProfile]
	cachedProfilesByNickname *cachestorage.CacheStorage[string, shikitypes.UserProfile]
	cachedTopics             *cachestorage.CacheStorage[uint, shikitypes.Topic]
	cachedImagesLink         *cachestorage.CacheStorage[uint, string]
}

func NewCacher() *Cacher {
	return &Cacher{
		cachedComments:           cachestorage.NewCacheStorage[uint, shikitypes.Comment](),
		cachedProfiles:           cachestorage.NewCacheStorage[uint, shikitypes.UserProfile](),
		cachedProfilesByNickname: cachestorage.NewCacheStorage[string, shikitypes.UserProfile](),
		cachedTopics:             cachestorage.NewCacheStorage[uint, shikitypes.Topic](),
		cachedImagesLink:         cachestorage.NewCacheStorage[uint, string](),
	}
}

func (c *Cacher) NumCached() uint {
	numCached := c.cachedComments.NumCached()
	numCached += c.cachedProfiles.NumCached()
	numCached += c.cachedProfilesByNickname.NumCached()
	numCached += c.cachedTopics.NumCached()
	numCached += c.cachedImagesLink.NumCached()
	return numCached
}

func (c *Cacher) Clear() uint {
	numDeleted := c.NumCached()
	return numDeleted
}

func (c *Cacher) SetComment(commentID uint, comment shikitypes.Comment) {
	c.cachedComments.Set(commentID, comment)
}

func (c *Cacher) GetComment(commentID uint) *shikitypes.Comment {
	return c.cachedComments.Get(commentID)
}

func (c *Cacher) SetProfile(profileID uint, profile shikitypes.UserProfile) {
	c.cachedProfiles.Set(profileID, profile)
}

func (c *Cacher) GetProfile(profileID uint) *shikitypes.UserProfile {
	return c.cachedProfiles.Get(profileID)
}

func (c *Cacher) SetProfileByNickname(nickname string, profile shikitypes.UserProfile) {
	c.cachedProfilesByNickname.Set(nickname, profile)
}

func (c *Cacher) GetProfileByNickname(nickname string) *shikitypes.UserProfile {
	return c.cachedProfilesByNickname.Get(nickname)
}

func (c *Cacher) SetTopic(topicID uint, topic shikitypes.Topic) {
	c.cachedTopics.Set(topicID, topic)
}

func (c *Cacher) GetTopic(topicID uint) *shikitypes.Topic {
	return c.cachedTopics.Get(topicID)
}

func (c *Cacher) SetImage(imageID uint, imageLink string) {
	c.cachedImagesLink.Set(imageID, imageLink)
}

func (c *Cacher) GetImage(imageID uint) *string {
	return c.cachedImagesLink.Get(imageID)
}
