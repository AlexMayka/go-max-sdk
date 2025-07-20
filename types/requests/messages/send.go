package messages

import "github.com/AlexMayka/go-max-sdk/types/models"

// Send представляет запрос на отправку сообщения (POST /messages)
type Send struct {
	// Если вы хотите отправить сообщение пользователю, укажите его ID (опционально)
	UserID *int64 `json:"-" query:"user_id,omitempty"`
	
	// Если сообщение отправляется в чат, укажите его ID (опционально)
	ChatID *int64 `json:"-" query:"chat_id,omitempty"`
	
	// Если false, сервер не будет генерировать превью для ссылок в тексте сообщения (опционально)
	DisableLinkPreview *bool `json:"-" query:"disable_link_preview,omitempty"`
	
	// Тело сообщения
	models.NewMessageBody
}