package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"gorm.io/gorm"
)

type Class struct {
	gorm.Model
	CreateByID     uint
	CreateBy       User
	TeacherID      uint
	Teacher        User
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

func ClassToAPI(c Class) types.Class {
	return types.Class{
		ID:             c.ID,
		CreatedByID:    c.CreateByID,
		CreatedBy:      UserToAPI(c.CreateBy),
		TeacherID:      c.TeacherID,
		Teacher:        UserToAPI(c.Teacher),
		Code:           c.Code,
		Name:           c.Name,
		CourseID:       c.CourseID,
		Course:         CourseToAPI(c.Course),
		Paralelo:       c.Paralelo,
		AcademicPeriod: c.AcademicPeriod,
		Description:    c.Description,
		ImgBackURL:     c.ImgBackURL,
		IsPublic:       c.IsPublic,
	}
}

func RegisterClass(classAPI types.Class) (id uint, err error) {
	class := Class{
		CreateByID:     classAPI.CreatedByID,
		TeacherID:      classAPI.TeacherID,
		Code:           classAPI.Code,
		Name:           classAPI.Name,
		CourseID:       classAPI.CourseID,
		Paralelo:       classAPI.Paralelo,
		AcademicPeriod: classAPI.AcademicPeriod,
		Description:    classAPI.Description,
		ImgBackURL:     classAPI.ImgBackURL,
		IsPublic:       classAPI.IsPublic,
	}
	result := db.DB.Create(&class)
	if result.Error != nil {
		return 0, result.Error
	}

	return class.ID, nil
}

func GetClassByID(id uint) (Class, error) {
	var class Class
	class.ID = id
	result := db.DB.Preload("CreateBy").Preload("Teacher").Preload("Course").First(&class)
	if result.Error != nil {
		return class, result.Error
	}
	return class, nil
}
