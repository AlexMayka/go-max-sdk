package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// RemoveMember представляет результат удаления участника из чата (ответ на DELETE /chats/{chatId}/members)
type RemoveMember struct {
	models.SuccessResponse
}