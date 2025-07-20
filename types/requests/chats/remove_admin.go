package chats

// RemoveAdmin представляет запрос на отмену прав администратора (DELETE /chats/{chatId}/members/admins/{userId})
type RemoveAdmin struct {
	// ID чата
	ChatID int64 `json:"-" path:"chatId"`
	
	// Идентификатор пользователя
	UserID int64 `json:"-" path:"userId"`
}