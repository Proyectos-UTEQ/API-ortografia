package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"fmt"

	"gorm.io/gorm"
)

type Module struct {
	gorm.Model
	CreatedByID      uint
	CreatedBy        User `gorm:"foreignKey:CreatedByID"`
	Title            string
	ShortDescription string
	TextRoot         string
	ImgBackURL       string
	Difficulty       Difficulty
	PointsToEarn     int
	Index            int
	IsPublic         bool
}

type Difficulty string

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)

func (Module) TableName() string {
	return "modules"
}

func RegisterModuleForTeacher(module *types.Module, userid uint) (*types.Module, error) {
	moduledb := Module{
		CreatedByID:      userid,
		Title:            module.Title,
		ShortDescription: module.ShortDescription,
		TextRoot:         module.TextRoot,
		ImgBackURL:       module.ImgBackURL,
		Difficulty:       Difficulty(module.Difficulty),
		PointsToEarn:     module.PointsToEarn,
		Index:            module.Index,
		IsPublic:         module.IsPublic,
	}

	// guardamos el modulos en la db
	result := db.DB.Create(&moduledb)
	if result.Error != nil {
		return nil, result.Error
	}

	// recuperamos el usuario de la db.

	result = db.DB.Preload("CreatedBy").First(&moduledb, moduledb.ID)
	if result.Error != nil {
		return nil, result.Error
	}

	fmt.Println(module)

	return &types.Module{
		ID:        moduledb.ID,
		CreatedAt: moduledb.CreatedAt.String(),
		UpdatedAt: moduledb.UpdatedAt.String(),
		CreateBy: types.UserAPI{
			ID:        moduledb.CreatedBy.ID,
			Email:     moduledb.CreatedBy.Email,
			FirstName: moduledb.CreatedBy.FirstName,
			LastName:  moduledb.CreatedBy.LastName,
			URLAvatar: moduledb.CreatedBy.URLAvatar,
		},
		Title:            moduledb.Title,
		ShortDescription: moduledb.ShortDescription,
		TextRoot:         moduledb.TextRoot,
		ImgBackURL:       moduledb.ImgBackURL,
		Difficulty:       string(moduledb.Difficulty),
		PointsToEarn:     moduledb.PointsToEarn,
		Index:            moduledb.Index,
		IsPublic:         moduledb.IsPublic,
	}, nil

}
