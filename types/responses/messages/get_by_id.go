package messages

import "github.com/AlexMayka/go-max-sdk/types/models"

// GetByID представляет сообщение по ID (ответ на GET /messages/{messageId})
type GetByID struct {
	models.Message
}