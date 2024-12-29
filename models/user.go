package models

import (
	gorm "gorm.io/gorm"
)

type User struct {
	gorm.Model
	TgID  int64 `gorm:"unique;notnull"`
	Level uint
}
