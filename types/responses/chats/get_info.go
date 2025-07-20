package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// GetInfo представляет информацию о чате (ответ на GET /chats/{chatId})
// Возвращает полную информацию о чате включая participants и pinned_message
type GetInfo struct {
	models.Chat
}