package messages

import "github.com/AlexMayka/go-max-sdk/types/models"

// Send представляет отправленное сообщение (ответ на POST /messages)
type Send struct {
	// Сообщение в чате
	Message models.Message `json:"message"`
}