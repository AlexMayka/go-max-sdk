package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// LeaveChat представляет результат удаления бота из чата (ответ на DELETE /chats/{chatId}/members/me)
type LeaveChat struct {
	models.SuccessResponse
}