package messages

import "github.com/AlexMayka/go-max-sdk/types/models"

// AnswerCallback представляет запрос ответа на callback (POST /answers)
type AnswerCallback struct {
	// Идентификатор кнопки, по которой пользователь кликнул (от 1 символа)
	CallbackID string `json:"-" query:"callback_id"`
	
	// Заполните это, если хотите изменить текущее сообщение (опционально)
	Message *models.NewMessageBody `json:"message,omitempty"`
	
	// Заполните это, если хотите просто отправить одноразовое уведомление пользователю (опционально)
	Notification *string `json:"notification,omitempty"`
}