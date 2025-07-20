package messages

import "github.com/AlexMayka/go-max-sdk/types/models"

// GetList представляет список сообщений (ответ на GET /messages)
type GetList struct {
	// Массив сообщений
	Messages []models.Message `json:"messages"`
}