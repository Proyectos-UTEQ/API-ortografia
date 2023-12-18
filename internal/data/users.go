package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string
	LastName     string
	Email        string `gorm:"uniqueIndex"`
	Password     string
	BirthDate    time.Time
	PointsEarned int
	Whatsapp     string
	Telegram     string
	URLAvatar    string
	Status       Status
	TypeUser     TypeUser
}

type Status string

const (
	Actived Status = "actived"
	Blocked Status = "blocked"
)

type TypeUser string

const (
	Admin   TypeUser = "admin"
	Student TypeUser = "student"
	Teacher TypeUser = "teacher"
)

func (User) TableName() string {
	return "users"
}

func Login(login types.Login) (*User, bool, error) {
	var user User
	result := db.DB.First(&user, "email = ?", login.Email)

	if result.Error != nil {
		return nil, false, result.Error
	}

	return &user, user.Password == login.Password, nil
}
