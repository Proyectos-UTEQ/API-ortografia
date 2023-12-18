package data

import "gorm.io/gorm"

type Module struct {
	gorm.Model
	CreatedByID      uint
	CreatedBy        User `gorm:"foreignKey:CreatedByID"`
	Title            string
	ShortDescription string
	TextRoot         string
	ImgBackURL       string
	Difficulty       Difficulty
	PointsToEarn     string
	Index            int
	IsPublic         bool
}

type Difficulty string

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)
