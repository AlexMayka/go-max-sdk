package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// Update представляет обновленную информацию о чате (ответ на PATCH /chats/{chatId})
type Update struct {
	models.Chat
}