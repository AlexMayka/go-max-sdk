package models

// BotCommand представляет команду, поддерживаемую ботом
type BotCommand struct {
	// Название команды (от 1 до 64 символов)
	Name string `json:"name"`

	// Описание команды (по желанию) (от 1 до 128 символов) (опционально)
	Description *string `json:"description,omitempty"`
}
