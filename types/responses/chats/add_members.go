package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// AddMembers представляет результат добавления участников в чат (ответ на POST /chats/{chatId}/members)
type AddMembers struct {
	models.SuccessResponse
}