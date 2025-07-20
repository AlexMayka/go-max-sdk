package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// GetMyMembership представляет информацию о членстве бота в чате (ответ на GET /chats/{chatId}/members/me)
type GetMyMembership struct {
	models.ChatMember
}