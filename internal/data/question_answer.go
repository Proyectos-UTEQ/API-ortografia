package data

import (
	"Proyectos-UTEQ/api-ortografia/pkg/types"

	"github.com/lib/pq"
)

// QuestionAnswer registra las opciones de la pregunta.
type Options struct {
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

func QuestionAnswerToAPI(questionAnswer Options) types.Options {
	return types.Options{
		SelectMode:     string(questionAnswer.SelectMode),
		TextOptions:    questionAnswer.TextOptions,
		TextToComplete: questionAnswer.TextToComplete,
		Hind:           questionAnswer.Hind,
	}
}
