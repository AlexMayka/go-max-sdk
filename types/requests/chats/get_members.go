package chats

// GetMembers представляет запрос участников чата (GET /chats/{chatId}/members)
type GetMembers struct {
	// ID чата
	ChatID int64 `json:"-" path:"chatId"`
	
	// Список ID пользователей, чье членство нужно получить. Когда этот параметр передан, параметры count и marker игнорируются (опционально)
	UserIDs []int64 `json:"-" query:"user_ids,omitempty"`
	
	// Указатель на следующую страницу данных (опционально)
	Marker *int64 `json:"-" query:"marker,omitempty"`
	
	// Количество участников, которых нужно вернуть (от 1 до 100, по умолчанию 20) (опционально)
	Count *int `json:"-" query:"count,omitempty"`
}