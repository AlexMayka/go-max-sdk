package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// GetList представляет ответ со списком чатов (ответ на GET /chats)
type GetList struct {
	// Список запрашиваемых чатов
	Chats []models.Chat `json:"chats"`
	
	// Указатель на следующую страницу запрашиваемых чатов (может быть null)
	Marker *int64 `json:"marker,omitempty"`
}