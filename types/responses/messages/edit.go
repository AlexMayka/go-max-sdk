package messages

import "github.com/AlexMayka/go-max-sdk/types/models"

// Edit представляет результат редактирования сообщения (ответ на PUT /messages/{messageId})
type Edit struct {
	// Обновленное сообщение
	Message models.Message `json:"message"`
}