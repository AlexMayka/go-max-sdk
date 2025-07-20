package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// SendAction представляет результат отправки действия в чат (ответ на POST /chats/{chatId}/actions)
type SendAction struct {
	models.SuccessResponse
}