package data

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	ModuleID         uint
	Module           Module
	QuestionnaireID  uint
	Questionnaire    Questionnaire
	TextRoot         string
	Difficulty       int
	TypeQuestion     TypeQuestion
	QuestionAnswerID uint
	QuestionAnswer   QuestionAnswer
	CorrectAnswerID  uint
	CorrectAnswer    Answer
}

type TypeQuestion string

const (
	TrueFalse       TypeQuestion = "true_false"
	MultiChoiceText TypeQuestion = "multi_choice_text"
	MultiChoiceABC  TypeQuestion = "multi_choice_abc"
	CompleteWord    TypeQuestion = "complete_word"
	OrderWord       TypeQuestion = "order_word"
)
