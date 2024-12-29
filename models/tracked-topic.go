package models

import "gorm.io/gorm"

type TrackedTopic struct {
	gorm.Model
	UserID        uint  `gorm:"notnull"`
	User          *User `gorm:"foreignKey:UserID"`
	TopicID       uint  `gorm:"notnull"`
	LastCommentID uint  `gorm:"notnull"`
}
