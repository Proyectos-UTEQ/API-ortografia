package types

import "fmt"

type Class struct {
	ID             uint     `json:"id"`
	CreatedByID    uint     `json:"created_by_id"`
	CreatedBy      *UserAPI `json:"created_by"`
	TeacherID      uint     `json:"teacher_id"`
	Teacher        *UserAPI `json:"teacher"`
	Code           string   `json:"code"`
	Name           string   `json:"name" validate:"required"`
	CourseID       uint     `json:"course_id" validate:"required"`
	Course         Course   `json:"course"`
	Paralelo       string   `json:"paralelo" validate:"required"`
	AcademicPeriod string   `json:"academic_period" validate:"required"`
	Description    string   `json:"description" validate:"required"`
	ImgBackURL     string   `json:"img_back_url"`
	IsPublic       bool     `json:"is_public" validate:"required"`
}

func (c *Class) ValidateNewClass() error {
	if c.TeacherID == 0 {
		return fmt.Errorf("teacher_id is required")
	}

	if c.Name == "" {
		return fmt.Errorf("name is required")
	}

	if c.CourseID == 0 {
		return fmt.Errorf("course_id is required")
	}

	if c.Paralelo == "" {
		return fmt.Errorf("paralelo is required")
	}

	if c.AcademicPeriod == "" {
		return fmt.Errorf("academic_period is required")
	}

	if c.Description == "" {
		return fmt.Errorf("description is required")
	}

	return nil
}

type Course struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
