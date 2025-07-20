package chats

// RemoveMember представляет запрос на удаление участника из чата (DELETE /chats/{chatId}/members)
type RemoveMember struct {
	// ID чата
	ChatID int64 `json:"-" path:"chatId"`
	
	// ID пользователя, которого нужно удалить из чата
	UserID int64 `json:"-" query:"user_id"`
	
	// Если установлено в true, пользователь будет заблокирован в чате. Применяется только для чатов с публичной или приватной ссылкой (опционально)
	Block *bool `json:"-" query:"block,omitempty"`
}