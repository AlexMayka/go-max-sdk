package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// GetAdmins представляет список администраторов чата (ответ на GET /chats/{chatId}/members/admins)
type GetAdmins struct {
	// Список участников чата с информацией о времени последней активности
	Members []models.ChatMember `json:"members"`
	
	// Указатель на следующую страницу данных (опционально)
	Marker *int64 `json:"marker,omitempty"`
}