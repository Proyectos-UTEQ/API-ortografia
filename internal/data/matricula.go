package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"fmt"
	"gorm.io/gorm"
)

type Matricula struct {
	gorm.Model
	UserID  uint
	User    User
	ClassID uint
	Class   Class
}

// EnrollUser registra un usuario en una clase.
func EnrollUser(userID uint, code string) (uint, error) {

	// recuperar el id de la clase.
	classID, err := GetClassIDByCode(code)
	if err != nil {
		return 0, fmt.Errorf("error al recuperar el id de la clase")
	}

	// validar si el usuario ya está matriculado en la clase.
	var matricula Matricula
	result := db.DB.Where("user_id = ? AND class_id = ?", userID, classID).First(&matricula)
	if result.Error == nil {
		return 0, fmt.Errorf("usuario ya matriculado")
	}

	// registrar el usuario en la clase.
	matricula.UserID = userID
	matricula.ClassID = classID
	err = db.DB.Create(&matricula).Error
	if err != nil {
		return 0, fmt.Errorf("error al registrar el usuario")
	}
	return matricula.ID, nil
}

// GetClassIDByCode recupera el id de la clase a partir del code. ;)
func GetClassIDByCode(code string) (uint, error) {

	var class Class
	result := db.DB.Where("code = ?", code).Select("id").First(&class)
	if result.Error != nil {
		return 0, result.Error
	}
	return class.ID, nil
}
