package data

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// QuestionAnswer registra las opciones de la pregunta.
type QuestionAnswer struct {
	gorm.Model
	SelectMode     SelectMode
	TextOpcions    pq.StringArray `gorm:"type:varchar(200)[]"`
	TextToComplete string
	Hind           string
}

type SelectMode string

const (
	SelectModeSingle   SelectMode = "single"
	SelectModeMultiple SelectMode = "multiple"
)
