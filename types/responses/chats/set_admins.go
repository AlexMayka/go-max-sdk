package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// SetAdmins представляет результат назначения администраторов (ответ на POST /chats/{chatId}/members/admins)
type SetAdmins struct {
	models.SuccessResponse
}