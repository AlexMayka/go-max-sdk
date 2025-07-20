package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// UnpinMessage представляет результат удаления закрепленного сообщения (ответ на DELETE /chats/{chatId}/pin)
type UnpinMessage struct {
	models.SuccessResponse
}