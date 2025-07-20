package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// Delete представляет результат удаления чата (ответ на DELETE /chats/{chatId})
type Delete struct {
	models.SuccessResponse
}