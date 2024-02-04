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
	Archived       bool
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
		Archived:       c.Archived,
	}
}

func ClassesToAPI(classes []Class) []types.Class {
	var classesAPI []types.Class
	for _, c := range classes {
		classesAPI = append(classesAPI, ClassToAPI(c))
	}
	return classesAPI
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
		Archived:       classAPI.Archived,
	}
	result := db.DB.Create(&class)
	if result.Error != nil {
		return 0, result.Error
	}

	return class.ID, nil
}

func UpdateClassByID(classAPI types.Class) error {

	// Actualizamos la clase.
	var class Class
	class.ID = classAPI.ID
	result := db.DB.Model(&class).Select("teacher_id", "name", "course_id", "paralelo", "academic_period", "description", "img_back_url", "archived").Updates(Class{
		TeacherID:      classAPI.TeacherID,
		Name:           classAPI.Name,
		CourseID:       classAPI.CourseID,
		Paralelo:       classAPI.Paralelo,
		AcademicPeriod: classAPI.AcademicPeriod,
		Description:    classAPI.Description,
		ImgBackURL:     classAPI.ImgBackURL,
		Archived:       classAPI.Archived,
	})

	if result.Error != nil {
		return result.Error
	}
	return nil
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

func GetClassesByTeacherID(teacherID uint) ([]Class, error) {
	var classes []Class
	result := db.DB.Preload("CreateBy").Preload("Teacher").Preload("Course").Where("teacher_id = ? and archived = false", teacherID).Find(&classes)
	if result.Error != nil {
		return classes, result.Error
	}
	return classes, nil
}

func GetClassesArchivedByTeacherID(teacherID uint) ([]Class, error) {
	var classes []Class
	result := db.DB.Preload("CreateBy").Preload("Teacher").Preload("Course").Where("teacher_id = ? and archived = true", teacherID).Find(&classes)
	if result.Error != nil {
		return classes, result.Error
	}
	return classes, nil
}

func ArchiveClass(id uint) error {
	var class Class
	class.ID = id
	result := db.DB.Model(&class).Update("archived", true)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
