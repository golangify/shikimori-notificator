package models

import "gorm.io/gorm"

type TrackedProfile struct {
	gorm.Model
	UserID              uint  `gorm:"notnull"` // models.User.ID
	User                *User `gorm:"foreignKey:UserID"`
	ProfileID           uint  `gorm:"notnull"`
	LastCommentID       uint  `gorm:"notnull"`
	TrackPosting        bool  `gorm:"default:true"`
	LastPostedCommentID uint
}
