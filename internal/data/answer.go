package data

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Answer struct {
	gorm.Model
	TrueOrFalse    bool
	TextOpcions    pq.StringArray `gorm:"type:varchar(200)[]"`
	TextToComplete pq.StringArray `gorm:"type:varchar(200)[]"`
}
