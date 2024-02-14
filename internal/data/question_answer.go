package data

import (
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"math/rand"
	"time"

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
	rand.Seed(time.Now().UnixNano())
	shuffledArray := make([]string, len(questionAnswer.TextOptions))
	copy(shuffledArray, questionAnswer.TextOptions)
	rand.Shuffle(len(shuffledArray), func(i, j int) {
		shuffledArray[i], shuffledArray[j] = shuffledArray[j], shuffledArray[i]
	})
	return types.Options{
		SelectMode:     string(questionAnswer.SelectMode),
		TextOptions:    shuffledArray,
		TextToComplete: questionAnswer.TextToComplete,
		Hind:           questionAnswer.Hind,
	}
}
