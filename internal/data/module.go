package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"Proyectos-UTEQ/api-ortografia/internal/utils"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"errors"
	"fmt"
	"math"
	"time"

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

//const (
//	Easy   Difficulty = "easy"
//	Medium Difficulty = "medium"
//	Hard   Difficulty = "hard"
//)

func (Module) TableName() string {
	return "modules"
}

// ModulesToAPI convierte las entidades de módulos a tipos de módulos para mostrar en la API REST xD
func ModulesToAPI(modules []Module) []types.Module {
	modulesApi := make([]types.Module, len(modules))
	for i, module := range modules {
		modulesApi[i] = ModuleToApi(module)
	}
	return modulesApi
}

// ModuleToApi convertimos un module data a un module type para la API REST.
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
		Difficulty:       DifficultyToFrontend(string(module.Difficulty)),
		PointsToEarn:     module.PointsToEarn,
		Index:            module.Index,
		IsPublic:         module.IsPublic,
	}
}

func DifficultyToFrontend(difficulty string) string {
	switch difficulty {
	case "easy":
		return "Fácil"
	case "medium":
		return "Medio"
	case "hard":
		return "Difícil"
	}

	return "Desconocido"
}

func RegisterModuleForTeacher(module *types.Module, userid uint) (types.Module, error) {

	moduleDB := Module{
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

	// guardamos el módulo en la db
	result := db.DB.Create(&moduleDB)
	if result.Error != nil {
		return types.Module{}, result.Error
	}

	// recuperamos el usuario de la db.
	result = db.DB.Preload("CreatedBy").First(&moduleDB, moduleDB.ID)
	if result.Error != nil {
		return types.Module{}, result.Error
	}

	return ModuleToApi(moduleDB), nil

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

// GetModulesForTeacher Se encarga de traer los módulos creados por el profesor.
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

	// establecemos la cantidad de items por página
	paginatedDetails.ItemsPerPage = len(modules)

	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) {
			if pgErr.Code == "42703" {
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
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) {
			if pgErr.Code == "42703" {
				return nil, nil, fmt.Errorf("columna inexistente: %s", paginated.Sort)
			}
		}
		return nil, nil, result.Error
	}

	return modules, &paginatedDetails, nil
}

// GetModule Se encarga de traer todos los módulos, sin importar quien los haya creado.
func GetModule(paginated *types.Paginated) (modules []Module, details types.PagintaedDetails, err error) {

	// cantidad total de módulos.
	db.DB.
		Table("modules").
		Where("title LIKE ?", "%"+paginated.Query+"%").
		Count(&details.TotalItems)

	// pagina actual y total de paginas.
	details.Page = paginated.Page
	details.TotalPage = int64(math.Ceil(float64(details.TotalItems) / float64(paginated.Limit)))

	// Recuperamos los módulos
	result := db.DB.
		Preload("CreatedBy").
		Where("title LIKE ?", "%"+paginated.Query+"%").
		Order(fmt.Sprintf("%s %s", paginated.Sort, paginated.Order)).
		Limit(paginated.Limit).
		Offset((paginated.Page - 1) * paginated.Limit).
		Find(&modules)

	// establecemos la cantidad de items por página
	details.ItemsPerPage = len(modules)

	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) {
			if pgErr.Code == "42703" {
				return nil, details, fmt.Errorf("columna inexistente: %s", paginated.Sort)
			}
		}
		return nil, details, result.Error
	}

	return modules, details, nil

}

// ModuleUserSub Recupera los datos necesário para un servicio en específico.
type ModuleUserSub struct {
	Module
	IsSubscribed bool
}

// ModuleUserSubToApi funcioné para convertir los módulos que agrega la parte si están suscritos o no,
// para utilizarlos en el frontend.
func ModuleUserSubToApi(modules []ModuleUserSub) []types.ModuleUser {

	modulesApi := make([]types.ModuleUser, len(modules))
	for i, module := range modules {
		modulesApi[i] = ModuleUserToApi(module)
	}
	return modulesApi

}

func ModuleUserToApi(module ModuleUserSub) types.ModuleUser {
	return types.ModuleUser{
		Module:       ModuleToApi(module.Module),
		IsSubscribed: module.IsSubscribed,
	}
}

// GetModuleWithUserSubscription Retorna todos los módulos y además tiene un campo para saber si el usuario está suscrito a ese módulo
func GetModuleWithUserSubscription(paginated *types.Paginated, userid uint) (moduleUser []ModuleUserSub, details types.PagintaedDetails, err error) {

	// cantidad total de módulos.
	db.DB.
		Table("modules").
		Where("title LIKE ? AND is_public = true", "%"+paginated.Query+"%").
		Count(&details.TotalItems)

	// pagina actual y total de paginas.
	details.Page = paginated.Page
	details.TotalPage = int64(math.Ceil(float64(details.TotalItems) / float64(paginated.Limit)))

	// Recuperamos todos los módulos, y en cada módulo revisamos si el usuario está suscrito.
	result := db.DB.
		Table("modules").
		Preload("CreatedBy").
		Select("modules.* ", "subscriptions.user_id IS NOT NULL as is_subscribed").
		Joins("LEFT JOIN subscriptions ON subscriptions.module_id = modules.id").
		Where("title LIKE ? AND is_public = true", "%"+paginated.Query+"%").
		Where("subscriptions.user_id = ? or subscriptions.user_id is null ", userid). // where s.user_id = 3 or s.user_id is null
		Order(fmt.Sprintf("%s %s", paginated.Sort, paginated.Order)).
		Limit(paginated.Limit).
		Offset((paginated.Page - 1) * paginated.Limit).
		Find(&moduleUser)

	// establecemos la cantidad de items por página
	details.ItemsPerPage = len(moduleUser)

	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) {
			if pgErr.Code == "42703" {
				return nil, details, fmt.Errorf("columna inexistente: %s", paginated.Sort)
			}
		}
		return nil, details, result.Error
	}

	return moduleUser, details, nil

}

// GetStudentsByModule StudentForModule recupera los estudiantes de un modulo
func GetStudentsByModule(moduleID uint) ([]User, error) {
	var users []User
	result := db.DB.
		Table("users").
		Joins("JOIN subscriptions ON subscriptions.user_id = users.id").
		Where("subscriptions.module_id = ?", moduleID).
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

// StudentPointsList Recupera una lista de estudiantes con la suma de puntajes que han obtenido en los módulos.
func StudentPointsList(start, end time.Time, limit int) ([]types.PointsUserForModule, error) {
	var pointsList []types.PointsUserForModule
	result := db.DB.Table("test_modules").
		Select("user_id, sum(test_modules.qualification) as points, users.first_name, users.last_name, users.url_avatar").
		Joins("JOIN users ON users.id = test_modules.user_id").
		Where("test_modules.created_at BETWEEN ? AND ?", start, end).
		Limit(limit).
		Group("user_id, users.first_name, users.last_name, users.url_avatar").Find(&pointsList)

	if result.Error != nil {
		return nil, result.Error
	}

	for i := range pointsList {
		if pointsList[i].URLAvatar == "" {
			pointsList[i].URLAvatar = fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=5952A2&color=fff&size=128", pointsList[i].FirstName+pointsList[i].LastName)
		}
	}
	return pointsList, nil
}
