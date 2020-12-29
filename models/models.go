package models

import "gorm.io/gorm"

type Installation struct {
	gorm.Model
	DeviceToken string `gorm:"size:128;uniqueIndex"`
	Server string `gorm:"size:32;index"`
	Locale string `gorm:"size:32;index"`
	ClientPreferences string `gorm:"size:256;index"`
	ClientVersion string `gorm:"size:32;index"`
}
