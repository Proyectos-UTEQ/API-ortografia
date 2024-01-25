package data

import "gorm.io/gorm"

type Questionnaire struct {
	gorm.Model
	ClassID uint
	Class   Class
	Title   string
}
