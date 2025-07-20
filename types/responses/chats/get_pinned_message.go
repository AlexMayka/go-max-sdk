package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// GetPinnedMessage представляет закрепленное сообщение (ответ на GET /chats/{chatId}/pin)
type GetPinnedMessage struct {
	// Закрепленное сообщение. Может быть null, если в чате нет закрепленного сообщения (опционально)
	Message *models.Message `json:"message,omitempty"`
}