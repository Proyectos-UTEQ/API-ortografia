package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"errors"

	"gorm.io/gorm"
)

type ResetPassword struct {
	gorm.Model
	UserID uint
	User   User
	Email  string
	Token  string
	Used   bool
}

func (ResetPassword) TableName() string {
	return "reset_password"
}

// registramos un nuevo reset password
// basicamente se registra que se genero un nuevo jwt.
func SaveResetPassword(id uint, email string, token string) error {
	resetPassword := ResetPassword{
		UserID: id,
		Email:  email,
		Token:  token,
		Used:   false,
	}

	result := db.DB.Create(&resetPassword)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func TokenIsUsed(token string) (bool, ResetPassword, error) {
	var resetPassword ResetPassword
	result := db.DB.Model(&ResetPassword{}).Where("token = ?", token).First(&resetPassword)
	if result.Error != nil {
		return false, resetPassword, result.Error
	}
	if resetPassword.Used {
		return true, resetPassword, errors.New("token already used")
	}
	return false, resetPassword, nil
}

func SetTokenUsed(token string) error {
	result := db.DB.Model(&ResetPassword{}).Where("token = ?", token).Update("used", true)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
