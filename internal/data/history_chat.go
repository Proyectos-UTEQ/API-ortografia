package data

import "gorm.io/gorm"

type HistoryChat struct {
	gorm.Model
	ChatIssueID uint
	ChatIssue   ChatIssue
	Message     string
	IsIA        bool
}
