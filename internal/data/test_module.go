package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"Proyectos-UTEQ/api-ortografia/internal/utils"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"time"

	"gorm.io/gorm"
)

type TestModule struct {
	gorm.Model
	UserID        uint
	User          User
	ModuleID      uint
	Module        Module
	Started       *time.Time
	Finished      *time.Time
	Qualification float32
}

func TestModuleToAPI(testModule TestModule) types.TestModule {
	return types.TestModule{
		ID:                        testModule.ID,
		CreatedAt:                 utils.GetFullDate(testModule.CreatedAt),
		ModuleID:                  testModule.Module.ID,
		Module:                    ModuleToApi(testModule.Module),
		Started:                   utils.GetFullDateOrNull(testModule.Started),
		Finished:                  utils.GetFullDateOrNull(testModule.Finished),
		Qualification:             testModule.Qualification,
		TestModuleQuestionAnswers: nil,
	}
}

// TestsModuleToAPI convierte una lista de testModule para el frontend.
func TestsModuleToAPI(testsModules []TestModule) []types.TestModule {
	var testModulesAPI []types.TestModule
	for _, v := range testsModules {
		testModulesAPI = append(testModulesAPI, TestModuleToAPI(v))
	}
	return testModulesAPI
}

func GenerateTestForStudent(userid uint, moduleID uint) (testId uint, err error) {

	// crear el objeto test Module

	tx := db.DB.Begin()
	now := time.Now()
	test := TestModule{
		UserID:        userid,
		ModuleID:      moduleID,
		Started:       &now,
		Finished:      nil,
		Qualification: 0,
	}

	// lo registramos en la base de datos.
	result := tx.Create(&test)
	if result.Error != nil {
		tx.Rollback()
		return 0, result.Error
	}

	// Seleccionamos 10 preguntas aleatorias del módulo.
	questions, err := GenerateQuestions(moduleID, 10)
	if err != nil {
		return 0, result.Error
	}

	// Asignamos las preguntas asignadas a la respuesta del usuario, y el puntaje.
	// En la respuesta del usuario se dejará blanca la respuesta del usuario.
	answerUser := make([]AnswerUser, 0)
	for i := range questions {
		answerUser = append(answerUser, AnswerUser{
			TestModuleID: test.ID,
			QuestionID:   questions[i].ID,
			Answer: Answer{
				TrueOrFalse:    false,
				TextOptions:    []string{""},
				TextToComplete: []string{""},
			},
			Score:       0,
			IsCorrect:   false,
			Responded:   false,
			Feedback:    "",
			ChatIssueID: nil,
		})
	}

	result = tx.Create(&answerUser)
	if result.Error != nil {
		tx.Rollback()
		return 0, result.Error
	}

	tx.Commit()

	return test.ID, nil
}

func TestByID(testid uint) (types.TestModule, error) {

	var test TestModule

	result := db.DB.Preload("Module.CreatedBy").Where("ID = ?", testid).Find(&test)
	if result.Error != nil {
		return types.TestModule{}, result.Error
	}

	// recuperamos las preguntas.
	var answerUser []AnswerUser
	result = db.DB.Preload("Answer").Preload("Question").Order("responded desc").Where("test_module_id = ?", test.ID).Find(&answerUser)
	if result.Error != nil {
		return types.TestModule{}, result.Error
	}

	responseModuleTest := types.TestModule{
		ID:            test.ID,
		CreatedAt:     test.CreatedAt.Format("02/01/2006 15:04:05"),
		ModuleID:      test.ModuleID,
		Module:        ModuleToApi(test.Module),
		Started:       utils.GetFullDateOrNull(test.Started),
		Finished:      utils.GetFullDateOrNull(test.Finished),
		Qualification: test.Qualification,
	}

	// recuperamos las respuestas del usuario.
	for i := range answerUser {
		questionAPI := QuestionToAPI(answerUser[i].Question)
		questionAPI.CorrectAnswerID = nil
		questionAPI.CorrectAnswer = nil
		responseModuleTest.TestModuleQuestionAnswers = append(responseModuleTest.TestModuleQuestionAnswers, types.TestModuleQuestionAnswer{
			Question:   questionAPI,
			AnswerUser: AnswerUserToAPI(answerUser[i]),
		})
	}

	return responseModuleTest, nil
}

func FinishTest(testid uint) (types.FinishTest, error) {
	// recuperar el test con las respuesta del usuario.
	tx := db.DB.Begin()

	var test TestModule
	tx.Where("ID = ?", testid).Find(&test)

	// Recupero todas las respuestas del usuario.
	var answersUser []AnswerUser
	tx.Where("test_module_id = ?", test.ID).Find(&answersUser)

	// calculo la calificacion.
	test.Qualification = 0
	for i := range answersUser {
		test.Qualification += float32(answersUser[i].Score)
	}
	now := time.Now()
	test.Finished = &now

	// actualizo el test.
	result := tx.Model(&TestModule{}).Select("qualification", "finished").Where("ID = ?", test.ID).Updates(test)
	if result.Error != nil {
		tx.Rollback()
		return types.FinishTest{}, result.Error
	}

	tx.Commit()

	return types.FinishTest{
		Finish:        test.Finished.Format("02/01/2006 15:04:05"),
		Qualification: test.Qualification,
		TestID:        test.ID,
	}, nil
}

func GetMyTest(userId, moduleId uint) ([]TestModule, error) {
	var testsModule []TestModule

	// Recuperamos los datos de la db.
	result := db.DB.
		Where("user_id = ? and module_id = ?", userId, moduleId).
		Preload("User").Preload("Module.CreatedBy").Find(&testsModule)

	if result.Error != nil {
		return nil, result.Error
	}

	return testsModule, nil
}
