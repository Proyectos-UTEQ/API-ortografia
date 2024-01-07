package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Module struct {
	gorm.Model
	CreatedByID      uint
	CreatedBy        User   `gorm:"foreignKey:CreatedByID"`
	Code             string `gorm:"uniqueIndex"`
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
		Code:             uuid.NewString(),
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
		Code:             moduledb.Code,
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

// Se encarga de traer los modulos creado por el profesor.
func GetModulesForTeacher(paginated *types.Paginated, userid uint) ([]Module, *types.PagintaedDetails, error) {

	var modules []Module
	var paginatedDetails types.PagintaedDetails

	// Calcular los detalles de la paginación.
	db.DB.
		Table("modules").
		Where("title LIKE ?", "%"+paginated.Query+"%").
		Where("created_by_id = ?", userid).Count(&paginatedDetails.TotalItems)
	paginatedDetails.Page = paginated.Page
	paginatedDetails.TotalPage = int64(math.Ceil(float64(paginatedDetails.TotalItems) / float64(paginated.Limit)))

	result := db.DB.
		Preload("CreatedBy").
		Where("title LIKE ?", "%"+paginated.Query+"%").
		Where("created_by_id = ?", userid).
		Order(fmt.Sprintf("%s %s", paginated.Sort, paginated.Order)).
		Limit(paginated.Limit).
		Offset((paginated.Page - 1) * paginated.Limit).
		Find(&modules)

	// seteamos la cantidad de items por pagina
	paginatedDetails.ItemsPerPage = len(modules)

	if result.Error != nil {
		if pgerr, ok := result.Error.(*pgconn.PgError); ok {
			if pgerr.Code == "42703" {
				return nil, nil, fmt.Errorf("columna inexistente: %s", paginated.Sort)
			}
		}
		return nil, nil, result.Error
	}
	return modules, &paginatedDetails, nil

}

// Se encarga de traer todos los modulos, sin importar quien los haya creado.
func GetModule(paginated *types.Paginated) (modules []Module, details types.PagintaedDetails, err error) {

	// cantidad total de modulos.
	db.DB.
		Table("modules").
		Where("title LIKE ?", "%"+paginated.Query+"%").
		Count(&details.TotalItems)

	// pagina actual y total de paginas.
	details.Page = paginated.Page
	details.TotalPage = int64(math.Ceil(float64(details.TotalItems) / float64(paginated.Limit)))

	// Recuperamos los modulos
	result := db.DB.
		Preload("CreatedBy").
		Where("title LIKE ?", "%"+paginated.Query+"%").
		Order(fmt.Sprintf("%s %s", paginated.Sort, paginated.Order)).
		Limit(paginated.Limit).
		Offset((paginated.Page - 1) * paginated.Limit).
		Find(&modules)

	// seteamos la cantidad de items por pagina
	details.ItemsPerPage = len(modules)

	if result.Error != nil {
		if pgerr, ok := result.Error.(*pgconn.PgError); ok {
			if pgerr.Code == "42703" {
				return nil, details, fmt.Errorf("columna inexistente: %s", paginated.Sort)
			}
		}
		return nil, details, result.Error
	}

	return modules, details, nil

}
