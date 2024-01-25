package data

import "gorm.io/gorm"

type Matricula struct {
	gorm.Model
	UserID  uint
	User    User
	ClassID uint
	Class   Class
}
