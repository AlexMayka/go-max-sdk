package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// PinMessage представляет результат закрепления сообщения (ответ на PUT /chats/{chatId}/pin)
type PinMessage struct {
	models.SuccessResponse
}