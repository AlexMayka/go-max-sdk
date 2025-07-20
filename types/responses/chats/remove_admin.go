package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// RemoveAdmin представляет результат отмены прав администратора (ответ на DELETE /chats/{chatId}/members/admins/{userId})
type RemoveAdmin struct {
	models.SuccessResponse
}