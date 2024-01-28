package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"Proyectos-UTEQ/api-ortografia/pkg/types"

	"gorm.io/gorm"
)

type AnswerUser struct {
	gorm.Model
	TestModuleID uint
	TestModule   TestModule // Relacion 1:1
	QuestionID   uint
	Question     Question // Pregunta en cual se basa la respuesta
	AnswerID     uint
	Answer       Answer // Respuesta del usuario
	Responded    bool   // Indica si el usuario respondio la pregunta
	Score        float32
	IsCorrect    bool
	Feedback     string
	ChatIssueID  *uint
	ChatIssue    ChatIssue
}

func AnswerUserToAPI(a AnswerUser) types.AnswerUser {
	return types.AnswerUser{
		AnswerUserID: a.ID,
		TestModuleID: a.TestModuleID,
		QuestionID:   a.QuestionID,
		Question:     nil,
		AnswerID:     a.AnswerID,
		Answer:       AnswerToAPI(&a.Answer),
		Responded:    a.Responded,
		Score:        a.Score,
		IsCorrect:    a.IsCorrect,
		Feedback:     a.Feedback,
	}
}

func GetAnswerUserByID(id uint) (AnswerUser, error) {
	var answerUser AnswerUser
	result := db.DB.Preload("Question.CorrectAnswer").First(&answerUser, id)
	return answerUser, result.Error
}

func (a *AnswerUser) UpdateAnswerUser() error {
	// se tiene que registrar el answer user y la respuesta en la otra tabla.
	tx := db.DB.Begin()

	result := tx.Model(&AnswerUser{}).Select("score", "is_correct", "responded", "feedback").Where("id = ?", a.ID).Updates(a)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// actualizamos la tabla answer
	result = tx.Model(&Answer{}).Select("true_or_false", "text_opcions", "text_to_complete").Where("id = ?", a.AnswerID).Updates(a.Answer)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil
}
