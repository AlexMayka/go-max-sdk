package models

// User представляет пользователя
type User struct {
	// ID пользователя
	UserID int64 `json:"user_id"`
	// Отображаемое имя пользователя
	FirstName string `json:"first_name"`
	// Отображаемая фамилия пользователя (может быть null)
	LastName string `json:"last_name"`
	// Устаревшее поле, скоро будет удалено (опционально)
	Name string `json:"name,omitempty"`
	// Уникальное публичное имя пользователя. Может быть null, если пользователь недоступен или имя не задано
	Username string `json:"username,omitempty"`
	// true, если пользователь является ботом
	IsBot bool `json:"is_bot"`
	// Время последней активности пользователя в MAX (Unix-время в миллисекундах). Может быть неактуальным, если пользователь отключил статус "онлайн" в настройках
	LastActivityTime int64 `json:"last_activity_time"`
}
