package messages

import "github.com/AlexMayka/go-max-sdk/types/models"

// Edit представляет запрос на редактирование сообщения (PUT /messages)
type Edit struct {
	// ID редактируемого сообщения (от 1 символа)
	MessageID string `json:"-" query:"message_id"`

	// Тело сообщения для редактирования
	models.NewMessageBody
}
