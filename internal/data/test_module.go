package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
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
	Started       time.Time
	Finished      time.Time
	Qualification float32
}

func GenerateTestForStudent(userid uint, moduleid uint) (testid uint, err error) {

	// crear el objeto test Module

	tx := db.DB.Begin()

	test := TestModule{
		UserID:        userid,
		ModuleID:      moduleid,
		Started:       time.Now(),
		Qualification: 0,
	}

	// lo registramos en la base de datos.
	result := tx.Create(&test)
	if result.Error != nil {
		tx.Rollback()
		return 0, result.Error
	}

	// Selecionamos 10 preguntas aleatorias del modulo.
	questions, err := GenerateQuestions(moduleid, 10)
	if err != nil {
		return 0, result.Error
	}

	// asignamos las preguntas asignadas a la respuesta del usuario, y el puntaje.
	// en la respuesta del usuario se dejara blanca la respuesta del usuario.
	answerUser := make([]AnswerUser, 0)
	for i := range questions {
		answerUser = append(answerUser, AnswerUser{
			TestModuleID: test.ID,
			QuestionID:   questions[i].ID,
			Answer: Answer{
				TrueOrFalse:    false,
				TextOpcions:    []string{""},
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
	result = db.DB.Preload("Answer").Preload("Question.QuestionAnswer").Order("responded desc").Where("test_module_id = ?", test.ID).Find(&answerUser)
	if result.Error != nil {
		return types.TestModule{}, result.Error
	}

	responseModuleTest := types.TestModule{
		ID:            test.ID,
		CreatedAt:     test.CreatedAt.Format("02/01/2006 15:04:05"),
		ModuleID:      test.ModuleID,
		Module:        ModuleToApi(test.Module),
		Started:       test.Started.Format("02/01/2006 15:04:05"),
		Finished:      test.Finished.Format("02/01/2006 15:04:05"),
		Qualification: test.Qualification,
	}

	// recuperamos las respuestas del usuario.
	for i := range answerUser {

		responseModuleTest.AnswerUser = append(responseModuleTest.AnswerUser, types.AnswerUser{
			AnswerUserID: answerUser[i].ID,
			TestModuleID: answerUser[i].TestModuleID,
			QuestionID:   answerUser[i].QuestionID,
			Question:     QuestionToAPI(answerUser[i].Question),
			AnswerID:     answerUser[i].AnswerID,
			Answer:       AnswerToAPI(answerUser[i].Answer),
			Score:        answerUser[i].Score,
			IsCorrect:    answerUser[i].IsCorrect,
			Responded:    answerUser[i].Responded,
			Feedback:     answerUser[i].Feedback,
			ChatIssueID:  answerUser[i].ChatIssueID,
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

	test.Finished = time.Now()

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
