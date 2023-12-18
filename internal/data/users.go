package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
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

func Login() (bool, error) {
	db.DB.Create(&User{
		FirstName:    "admin",
		LastName:     "admin",
		Email:        "roberto@gmail.com",
		Password:     "admin",
		BirthDate:    time.Now(),
		PointsEarned: 0,
		Whatsapp:     "",
		Telegram:     "",
		URLAvatar:    "",
	})
	return true, nil
}
