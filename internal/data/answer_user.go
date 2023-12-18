package data

import "gorm.io/gorm"

type AnswerUser struct {
	gorm.Model
	UserID       uint
	User         User
	QuestionID   uint
	Question     Question
	UserAnswerID uint
	//TODO: Falta agregar la relacion UserAnswer.
	Score       int
	IsCorrect   bool
	Feedback    string
	ChatIssueID uint
	//TODO: Falta agregar la relacion ChatIssue.
}
