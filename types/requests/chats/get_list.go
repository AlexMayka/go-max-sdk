package chats

// GetList представляет параметры запроса списка чатов (GET /chats)
type GetList struct {
	// Количество запрашиваемых чатов (от 1 до 100, по умолчанию 50) (опционально)
	Count *int `json:"-" query:"count,omitempty"`
	
	// Указатель на следующую страницу данных. Для первой страницы передайте null (опционально)
	Marker *int64 `json:"-" query:"marker,omitempty"`
}