package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Subscription representa una suscripción de un usuario a un modulo.
type Subscription struct {
	gorm.Model
	UserID   uint
	User     User
	ModuleID uint
	Module   Module
}

func (Subscription) TableName() string {
	return "subscriptions"
}

func RegisterSubscription(userID uint, code string) (Subscription, error) {

	// recuperar el id del module
	var module Module
	result := db.DB.First(&module, "code LIKE ?", fmt.Sprintf("%s%%", code))
	if result.RowsAffected == 0 {
		return Subscription{}, errors.New("el modulo no existe")
	}
	if result.Error != nil {
		return Subscription{}, result.Error

	}

	sub := Subscription{
		UserID:   userID,
		ModuleID: module.ID,
	}

	// Validar si este usuario ya se encuentra suscrito al módulo.
	result = db.DB.Where("user_id = ? AND module_id = ?", userID, module.ID).First(&sub)
	if result.Error == nil {
		return sub, errors.New("el usuario ya se encuentra suscrito al modulo")
	}

	result = db.DB.Create(&sub)
	if result.Error != nil {
		return sub, result.Error
	}

	return sub, nil

}
