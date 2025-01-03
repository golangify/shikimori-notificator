package models

import "gorm.io/gorm"

type TrackedProfile struct {
	gorm.Model
	UserID        uint  `gorm:"notnull"`
	User          *User `gorm:"foreignKey:UserID"`
	ProfileID     uint  `gorm:"notnull"`
	LastCommentID uint  `gorm:"notnull"`
}
