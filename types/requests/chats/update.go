package chats

import "github.com/AlexMayka/go-max-sdk/types/models"

// Update представляет запрос на изменение информации о чате (PATCH /chats/{chatId})
type Update struct {
	// ID чата
	ChatID int64 `json:"-" path:"chatId"`
	
	// Запрос на прикрепление изображения (все поля являются взаимоисключающими) (опционально)
	Icon *models.PhotoAttachmentRequestPayload `json:"icon,omitempty"`
	
	// Название чата (от 1 до 200 символов) (опционально)
	Title *string `json:"title,omitempty"`
	
	// ID сообщения для закрепления в чате. Чтобы удалить закреплённое сообщение, используйте метод unpin (опционально)
	Pin *string `json:"pin,omitempty"`
	
	// Если true, участники получат системное уведомление об изменении (по умолчанию true) (опционально)
	Notify *bool `json:"notify,omitempty"`
}