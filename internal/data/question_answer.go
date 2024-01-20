package data

import (
	"Proyectos-UTEQ/api-ortografia/pkg/types"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// QuestionAnswer registra las opciones de la pregunta.
type QuestionAnswer struct {
	gorm.Model
	SelectMode     SelectMode
	TextOptions    pq.StringArray `gorm:"type:varchar(200)[]"`
	TextToComplete string
	Hind           string
}

type SelectMode string

const (
	SelectModeSingle   SelectMode = "single"
	SelectModeMultiple SelectMode = "multiple"
)

func QuestionAnswerToAPI(questionAnswer QuestionAnswer) types.QuestionAnswer {
	return types.QuestionAnswer{
		ID:             questionAnswer.ID,
		SelectMode:     string(questionAnswer.SelectMode),
		TextOptions:    questionAnswer.TextOptions,
		TextToComplete: questionAnswer.TextToComplete,
		Hind:           questionAnswer.Hind,
	}
}
