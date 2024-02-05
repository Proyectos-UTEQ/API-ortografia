package interfaces

import "Proyectos-UTEQ/api-ortografia/pkg/types"

type QuestionIA interface {
	GenerateQuestion(typeQuestion string, text string) (*types.Question, error)
}
