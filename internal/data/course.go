package data

import (
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Name        string
	Description string
}

func CourseToAPI(c Course) types.Course {
	return types.Course{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
	}
}
