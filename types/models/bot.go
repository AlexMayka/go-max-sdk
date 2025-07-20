package models

// BotInfo представляет информацию о боте
type BotInfo struct {
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
	
	// Команды, поддерживаемые ботом (до 32 элементов) (опционально)
	Commands []*BotCommand `json:"commands,omitempty"`
}