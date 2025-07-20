package models

// Chat представляет чат
type Chat struct {
	// ID чата
	ChatID int64 `json:"chat_id"`

	// Тип чата (enum: "chat")
	Type ChatType `json:"type"`

	// Статус чата (enum: "active", "removed", "left", "closed")
	Status ChatStatus `json:"status"`

	// Отображаемое название чата. Может быть null для диалогов
	Title *string `json:"title,omitempty"`

	// Иконка чата
	Icon *Image `json:"icon,omitempty"`

	// Время последнего события в чате
	LastEventTime int64 `json:"last_event_time"`

	// Количество участников чата. Для диалогов всегда 2
	ParticipantsCount int32 `json:"participants_count"`

	// ID владельца чата (опционально)
	OwnerID *int64 `json:"owner_id,omitempty"`

	// Участники чата с временем последней активности. Может быть null, если запрашивается список чатов (опционально)
	Participants interface{} `json:"participants,omitempty"`

	// Доступен ли чат публично (для диалогов всегда false)
	IsPublic bool `json:"is_public"`

	// Ссылка на чат (опционально)
	Link *string `json:"link,omitempty"`

	// Описание чата
	Description *string `json:"description,omitempty"`

	// Данные о пользователе в диалоге (только для чатов типа "dialog") (опционально)
	DialogWithUser *UserWithPhoto `json:"dialog_with_user,omitempty"`

	// Количество сообщений в чате (доступно только для групповых чатов, недоступно для диалогов) (опционально)
	MessagesCount *int `json:"messages_count,omitempty"`

	// ID сообщения, содержащего кнопку, через которую был инициирован чат (опционально)
	ChatMessageID *string `json:"chat_message_id,omitempty"`

	// Закреплённое сообщение в чате (возвращается только при запросе конкретного чата) (опционально)
	PinnedMessage *Message `json:"pinned_message,omitempty"`
}

// ChatType представляет тип чата
type ChatType string

const (
	// ChatTypeChat групповой чат
	ChatTypeChat ChatType = "chat"
)

// ChatStatus представляет статус чата
type ChatStatus string

const (
	// ChatStatusActive бот является активным участником чата
	ChatStatusActive ChatStatus = "active"
	// ChatStatusRemoved бот был удалён из чата
	ChatStatusRemoved ChatStatus = "removed"
	// ChatStatusLeft бот покинул чат
	ChatStatusLeft ChatStatus = "left"
	// ChatStatusClosed чат был закрыт
	ChatStatusClosed ChatStatus = "closed"
)

// Image представляет изображение
type Image struct {
	// URL изображения
	URL string `json:"url"`
}

// UserWithPhoto представляет пользователя с фотографией
type UserWithPhoto struct {
	// ID пользователя
	UserID int64 `json:"user_id"`

	// Отображаемое имя пользователя
	FirstName string `json:"first_name"`

	// Отображаемая фамилия пользователя (может быть null)
	LastName *string `json:"last_name,omitempty"`

	// Устаревшее поле, скоро будет удалено (опционально)
	Name *string `json:"name,omitempty"`

	// Уникальное публичное имя пользователя. Может быть null, если пользователь недоступен или имя не задано
	Username *string `json:"username,omitempty"`

	// true, если пользователь является ботом
	IsBot bool `json:"is_bot"`

	// Время последней активности пользователя в MAX (Unix-время в миллисекундах). Может быть неактуальным, если пользователь отключил статус "онлайн" в настройках
	LastActivityTime int64 `json:"last_activity_time"`

	// Описание пользователя. Может быть null, если пользователь его не заполнил (до 16000 символов) (опционально)
	Description *string `json:"description,omitempty"`

	// URL аватара (опционально)
	AvatarURL *string `json:"avatar_url,omitempty"`

	// URL аватара большего размера (опционально)
	FullAvatarURL *string `json:"full_avatar_url,omitempty"`
}
