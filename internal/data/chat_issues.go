package data

import "gorm.io/gorm"

type ChatIssue struct {
	gorm.Model
	UserID uint
	User   User
	Issue  string
}
