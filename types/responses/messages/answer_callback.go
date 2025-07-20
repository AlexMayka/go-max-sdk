package messages

import "github.com/AlexMayka/go-max-sdk/types/models"

// AnswerCallback представляет результат ответа на callback (ответ на POST /answers)
type AnswerCallback struct {
	models.SuccessResponse
}