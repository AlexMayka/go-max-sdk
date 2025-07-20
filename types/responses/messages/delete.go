package messages

import "github.com/AlexMayka/go-max-sdk/types/models"

// Delete представляет результат удаления сообщения (ответ на DELETE /messages/{messageId})
type Delete struct {
	models.SuccessResponse
}