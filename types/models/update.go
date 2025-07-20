package models

// Update представляет различные типы событий, произошедших в чате
type Update struct {
	// Тип события
	UpdateType string `json:"update_type"`

	// Unix-время, когда произошло событие
	Timestamp int64 `json:"timestamp"`

	// Новое созданное сообщение
	Message *Message `json:"message,omitempty"`

	// Callback данные от inline-клавиатуры (только для message_callback)
	CallbackID *string `json:"callback_id,omitempty"`
	
	// Payload кнопки (только для message_callback)
	Payload *string `json:"payload,omitempty"`

	// Текущий язык пользователя в формате IETF BCP 47. Доступно только в диалогах (опционально)
	UserLocale *string `json:"user_locale,omitempty"`
}

// UpdateType представляет тип события
type UpdateType string

const (
	// UpdateTypeMessageCreated новое созданное сообщение
	UpdateTypeMessageCreated UpdateType = "message_created"
	
	// UpdateTypeMessageCallback callback от inline-клавиатуры
	UpdateTypeMessageCallback UpdateType = "message_callback"
	
	// UpdateTypeMessageEdited сообщение было отредактировано
	UpdateTypeMessageEdited UpdateType = "message_edited"
	
	// UpdateTypeMessageDeleted сообщение было удалено
	UpdateTypeMessageDeleted UpdateType = "message_deleted"
	
	// UpdateTypeChatMemberJoined участник присоединился к чату
	UpdateTypeChatMemberJoined UpdateType = "chat_member_joined"
	
	// UpdateTypeChatMemberLeft участник покинул чат
	UpdateTypeChatMemberLeft UpdateType = "chat_member_left"
)
