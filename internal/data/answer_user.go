package data

import "gorm.io/gorm"

type AnswerUser struct {
	gorm.Model
	UserID       uint
	User         User
	QuestionID   uint
	Question     Question
	UserAnswerID uint
	UserAnswer   Answer
	Score        int
	IsCorrect    bool
	Feedback     string
	ChatIssueID  uint
	ChatIssue    ChatIssue
}
