package data

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	Name        string
	Description string
}
