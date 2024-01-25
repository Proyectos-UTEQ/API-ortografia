package data

import (
	"Proyectos-UTEQ/api-ortografia/pkg/types"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Answer struct {
	gorm.Model
	TrueOrFalse    bool
	TextOpcions    pq.StringArray `gorm:"type:varchar(200)[]"`
	TextToComplete pq.StringArray `gorm:"type:varchar(200)[]"`
}

func AnswerToAPI(answer *Answer) *types.Answer {
	if answer == nil {
		return nil
	} else {
		return &types.Answer{
			ID:             answer.ID,
			TrueOrFalse:    answer.TrueOrFalse,
			TextOpcions:    answer.TextOpcions,
			TextToComplete: answer.TextToComplete,
		}
	}
}
