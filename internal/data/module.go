package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"Proyectos-UTEQ/api-ortografia/internal/utils"
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

// convierte las entidades de modulos a tipos de modulos para mostrar en la API REST xD
func ModulesToAPI(modules []Module, apphost string) []types.Module {
	// convertimos los modulos a types.modules
	modulesApi := make([]types.Module, len(modules))
	for i, module := range modules {
		modulesApi[i] = ModuleToApi(module)
	}
	return modulesApi
}

// convertimos un module data a un module type para la API REST.
func ModuleToApi(module Module) types.Module {
	return types.Module{
		ID:        module.ID,
		CreatedAt: utils.GetDate(module.CreatedAt),
		UpdatedAt: utils.GetDate(module.UpdatedAt),
		CreateBy: types.UserAPI{
			ID:                   module.CreatedBy.ID,
			FirstName:            module.CreatedBy.FirstName,
			LastName:             module.CreatedBy.LastName,
			Email:                module.CreatedBy.Email,
			URLAvatar:            module.CreatedBy.URLAvatar,
			Status:               string(module.CreatedBy.Status),
			TypeUser:             string(module.CreatedBy.TypeUser),
			PerfilUpdateRequired: module.CreatedBy.PerfilUpdateRequired,
		},
		Code:             module.Code,
		Title:            module.Title,
		ShortDescription: module.ShortDescription,
		TextRoot:         module.TextRoot,
		ImgBackURL:       module.ImgBackURL,
		Difficulty:       string(module.Difficulty),
		PointsToEarn:     module.PointsToEarn,
		Index:            module.Index,
		IsPublic:         module.IsPublic,
	}
}

func RegisterModuleForTeacher(module *types.Module, userid uint) (types.Module, error) {

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
		return types.Module{}, result.Error
	}

	// recuperamos el usuario de la db.
	result = db.DB.Preload("CreatedBy").First(&moduledb, moduledb.ID)
	if result.Error != nil {
		return types.Module{}, result.Error
	}

	return ModuleToApi(moduledb), nil

}

func UpdateModule(module *types.Module) (*Module, error) {
	data := map[string]interface{}{
		"title":             module.Title,
		"short_description": module.ShortDescription,
		"text_root":         module.TextRoot,
		"img_back_url":      module.ImgBackURL,
		"difficulty":        module.Difficulty,
		"points_to_earn":    module.PointsToEarn,
		"index":             module.Index,
		"is_public":         module.IsPublic,
	}

	result := db.DB.Model(&Module{}).Where("id = ?", module.ID).Updates(data)
	if result.Error != nil {
		return nil, result.Error
	}

	var moduleData Module
	result = db.DB.Preload("CreatedBy").First(&moduleData, module.ID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &moduleData, nil
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

func GetModuleForStudent(paginated *types.Paginated, userid uint) ([]Module, *types.PagintaedDetails, error) {
	var modules []Module
	var paginatedDetails types.PagintaedDetails

	db.DB.Model(&Module{}).
		Joins("JOIN subscriptions ON subscriptions.module_id = modules.id").
		Where("subscriptions.user_id = ?", userid).
		Count(&paginatedDetails.TotalItems)

	paginatedDetails.Page = paginated.Page
	paginatedDetails.TotalPage = int64(math.Ceil(float64(paginatedDetails.TotalItems) / float64(paginated.Limit)))

	result := db.DB.Model(&Module{}).
		Preload("CreatedBy").
		Joins("JOIN subscriptions ON subscriptions.module_id = modules.id").
		Where("subscriptions.user_id = ?", userid).
		Order(fmt.Sprintf("%s %s", paginated.Sort, paginated.Order)).
		Limit(paginated.Limit).
		Offset((paginated.Page - 1) * paginated.Limit).
		Find(&modules)

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

// Recupera los datos necesario para un servicio en especifico.
type ModuleUserSubcription struct {
	Module
	IsSubscribed bool
}

// funcione para convertir los modulos que agrega la parte si estan suscritos o no,
// para utilizarlos en el frontend.
func ModuleUserSubcriptionToApi(modules []ModuleUserSubcription) []types.ModuleUser {

	modulesApi := make([]types.ModuleUser, len(modules))
	for i, module := range modules {
		modulesApi[i] = ModuleUserToApi(module)
	}
	return modulesApi

}

func ModuleUserToApi(module ModuleUserSubcription) types.ModuleUser {
	return types.ModuleUser{
		Module:       ModuleToApi(module.Module),
		IsSubscribed: module.IsSubscribed,
	}
}

// Retorna todos los modulos y ademas tiene un campo para saber si el usuario esta suscrito a ese modulo
func GetModuleWithUserSubscription(paginated *types.Paginated, userid uint) (moduleUser []ModuleUserSubcription, details types.PagintaedDetails, err error) {

	// cantidad total de modulos.
	db.DB.
		Table("modules").
		Where("title LIKE ?", "%"+paginated.Query+"%").
		Count(&details.TotalItems)

	// pagina actual y total de paginas.
	details.Page = paginated.Page
	details.TotalPage = int64(math.Ceil(float64(details.TotalItems) / float64(paginated.Limit)))

	// Recuperamos todos los modulos, y en cada modulo revisamos si el usuario esta suscrito.
	result := db.DB.
		Table("modules").
		Preload("CreatedBy").
		Select("modules.* ", "subscriptions.user_id IS NOT NULL as is_subscribed").
		Joins("LEFT JOIN subscriptions ON subscriptions.module_id = modules.id").
		Where("title LIKE ?", "%"+paginated.Query+"%").
		Order(fmt.Sprintf("%s %s", paginated.Sort, paginated.Order)).
		Limit(paginated.Limit).
		Offset((paginated.Page - 1) * paginated.Limit).
		Find(&moduleUser)

	// seteamos la cantidad de items por pagina
	details.ItemsPerPage = len(moduleUser)

	if result.Error != nil {
		if pgerr, ok := result.Error.(*pgconn.PgError); ok {
			if pgerr.Code == "42703" {
				return nil, details, fmt.Errorf("columna inexistente: %s", paginated.Sort)
			}
		}
		return nil, details, result.Error
	}

	return moduleUser, details, nil

}

// StudentForModule recupera los estudiantes de un modulo
func GetStudentsByModule(moduleid uint) ([]User, error) {
	var users []User
	result := db.DB.
		Table("users").
		Joins("JOIN subscriptions ON subscriptions.user_id = users.id").
		Where("subscriptions.module_id = ?", moduleid).
		Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func ModuleByID(id uint) (*Module, error) {
	var module Module
	result := db.DB.Preload("CreatedBy").First(&module, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &module, nil
}
