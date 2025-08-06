package models

// Update представляет различные типы событий, произошедших в чате.
// Получается через polling (/updates) или webhook (/subscriptions).
// Каждое обновление имеет свой номер последовательности и тип события.
type Update struct {
	// Тип события
	UpdateType UpdateType `json:"update_type"`

	// Unix-время, когда произошло событие
	Timestamp int64 `json:"timestamp"`

	// Новое созданное сообщение
	Message *Message `json:"message,omitempty"`

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

	// UpdateTypeMessageRemoved сообщение было удалено
	UpdateTypeMessageRemoved UpdateType = "message_removed"

	// UpdateBotAdded бот был добавлен в чат
	UpdateBotAdded UpdateType = "bot_added"

	// UpdateBotRemoved бот был удален из чата
	UpdateBotRemoved UpdateType = "bot_removed"

	// UpdateUserAdded пользователь был добавлен в чат
	UpdateUserAdded UpdateType = "user_added"

	// UpdateUserRemoved пользователь был удален из чата или покинул чат
	UpdateUserRemoved UpdateType = "user_removed"

	// UpdateBotStarted пользователь нажал кнопку "Начать" в диалоге с ботом
	UpdateBotStarted UpdateType = "bot_started"

	// UpdateChatTitleChanged название чата было изменено
	UpdateChatTitleChanged UpdateType = "chat_title_changed"

	// UpdateMessageChatCreated чат был создан
	UpdateMessageChatCreated UpdateType = "message_chat_created"
)
