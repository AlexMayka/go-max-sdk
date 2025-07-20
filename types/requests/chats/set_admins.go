package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// SetAdmins представляет запрос на назначение администраторов (POST /chats/{chatId}/members/admins)
type SetAdmins struct {
	// ID чата
	ChatID int64 `json:"-" path:"chatId"`
	
	// Массив администраторов чата
	Admins []models.ChatAdmin `json:"admins"`
	
	// Указатель на следующую страницу данных (опционально)
	Marker *int64 `json:"marker,omitempty"`
}