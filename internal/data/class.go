package data

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	CreateByID     uint
	CreateBy       User
	Code           string
	Name           string
	CourseID       uint
	Course         Course
	Paralelo       string
	AcademicPeriod string
	Description    string
	ImgBackURL     string
	IsPublic       bool
}

func (Class) TableName() string {
	return "class"
}
