package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// GetByLink представляет информацию о чате по ссылке (ответ на GET /chats/{chatLink})
// Возвращает полную информацию о чате включая participants и pinned_message
type GetByLink struct {
	models.Chat
}