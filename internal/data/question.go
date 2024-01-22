package data

import (
	"Proyectos-UTEQ/api-ortografia/internal/db"
	"Proyectos-UTEQ/api-ortografia/pkg/types"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	ModuleID         *uint
	Module           Module
	QuestionnaireID  *uint
	Questionnaire    Questionnaire
	TextRoot         string
	Difficulty       int
	TypeQuestion     TypeQuestion
	QuestionAnswerID uint
	QuestionAnswer   QuestionAnswer
	CorrectAnswerID  uint
	CorrectAnswer    Answer
}

type TypeQuestion string

const (
	TrueFalse       TypeQuestion = "true_false"
	MultiChoiceText TypeQuestion = "multi_choice_text"
	MultiChoiceABC  TypeQuestion = "multi_choice_abc"
	CompleteWord    TypeQuestion = "complete_word"
	OrderWord       TypeQuestion = "order_word"
)

func QuestionToAPI(question Question) types.Question {
	return types.Question{
		ID:               question.ID,
		ModuleID:         question.ModuleID,
		QuestionnaireID:  question.QuestionnaireID,
		TextRoot:         question.TextRoot,
		Difficulty:       question.Difficulty,
		TypeQuestion:     string(question.TypeQuestion),
		QuestionAnswerID: question.QuestionAnswerID,
		QuestionAnswer:   QuestionAnswerToAPI(question.QuestionAnswer),
		CorrectAnswerID:  question.CorrectAnswerID,
		CorrectAnswer:    AnswerToAPI(question.CorrectAnswer),
	}
}

func QuestionListToAPI(questions []Question) []types.Question {
	questionList := make([]types.Question, 0)
	for _, question := range questions {
		questionList = append(questionList, QuestionToAPI(question))
	}
	return questionList
}

func RegisterQuestionForModule(questionAPI types.Question) (types.Question, error) {

	// convertimos los datos que nos envian a una question entidad
	question := Question{
		ModuleID:        questionAPI.ModuleID,
		QuestionnaireID: nil,
		TextRoot:        questionAPI.TextRoot,
		Difficulty:      questionAPI.Difficulty,
		TypeQuestion:    TypeQuestion(questionAPI.TypeQuestion),
		QuestionAnswer: QuestionAnswer{
			SelectMode:     SelectMode(questionAPI.QuestionAnswer.SelectMode),
			TextOptions:    pq.StringArray(questionAPI.QuestionAnswer.TextOptions),
			TextToComplete: questionAPI.QuestionAnswer.TextToComplete,
			Hind:           questionAPI.QuestionAnswer.Hind,
		},
		CorrectAnswer: Answer{
			TrueOrFalse:    questionAPI.CorrectAnswer.TrueOrFalse,
			TextOpcions:    pq.StringArray(questionAPI.CorrectAnswer.TextOpcions),
			TextToComplete: pq.StringArray(questionAPI.CorrectAnswer.TextToComplete),
		},
	}

	// Registramos en la base de datos.
	result := db.DB.Create(&question)
	if result.Error != nil {
		return QuestionToAPI(question), result.Error
	}

	return QuestionToAPI(question), nil
}

// Recupearmos todas las preguntas que pertenezcan al modulo
func GetQuestionsForModule(moduleID uint) ([]types.Question, error) {
	questions := []Question{}
	result := db.DB.Where("module_id = ?", moduleID).Preload("QuestionAnswer").Preload("CorrectAnswer").Find(&questions)
	if result.Error != nil {
		return nil, result.Error
	}
	return QuestionListToAPI(questions), nil
}
