package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
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

func Login(login types.Login) (*types.UserAPI, bool, error) {
	var user User
	result := db.DB.First(&user, "email = ?", login.Email)

	// Controlar el error de record not found.
	if result.Error != nil {
		return nil, false, errors.New("las credenciales son incorrectas")
	}

	// Comparar las contraseñas con un hash.
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		return nil, false, errors.New("las credenciales son incorrectas")
	}

	// Convertir a un usuario api
	userAPI := &types.UserAPI{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Password:     "",
		BirthDate:    user.BirthDate.String(),
		PointsEarned: user.PointsEarned,
		Whatsapp:     user.Whatsapp,
		Telegram:     user.Telegram,
		URLAvatar:    user.URLAvatar,
		Status:       string(user.Status),
		TypeUser:     string(user.TypeUser),
	}

	return userAPI, true, nil
}

func Register(userAPI *types.UserAPI) error {

	// crear un hash apartir de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userAPI.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// rellenamos los datos con la entidad.
	user := User{
		FirstName: userAPI.FirstName,
		Email:     userAPI.Email,
		Password:  string(hashedPassword),
		Status:    Actived,
		TypeUser:  TypeUser(userAPI.TypeUser),
	}

	result := db.DB.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	userAPI.ID = user.ID

	return nil
}

func ExisteEmail(email string) (bool, types.UserAPI) {
	var user User
	result := db.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return false, types.UserAPI{}
	}
	return true, types.UserAPI{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		BirthDate:    user.BirthDate.String(),
		PointsEarned: user.PointsEarned,
		Whatsapp:     user.Whatsapp,
		Telegram:     user.Telegram,
		URLAvatar:    user.URLAvatar,
		Status:       string(user.Status),
		TypeUser:     string(user.TypeUser),
	}
}

func UpdatePassword(userid uint, newPassword string) error {
	// hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// actualizar la contraseña en la base de datos.
	result := db.DB.Model(&User{}).Where("id = ?", userid).Update("password", string(hashedPassword))
	fmt.Println("Rows affected: ", result.RowsAffected)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
