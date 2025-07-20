package messages

// GetList представляет запрос сообщений (GET /messages)
type GetList struct {
	// ID чата, чтобы получить сообщения из определённого чата (опционально)
	ChatID *int64 `json:"-" query:"chat_id,omitempty"`
	
	// Список ID сообщений, которые нужно получить (через запятую) (опционально)
	MessageIDs []string `json:"-" query:"message_ids,omitempty"`
	
	// Время начала для запрашиваемых сообщений (в формате Unix timestamp) (опционально)
	From *int64 `json:"-" query:"from,omitempty"`
	
	// Время окончания для запрашиваемых сообщений (в формате Unix timestamp) (опционально)
	To *int64 `json:"-" query:"to,omitempty"`
	
	// Максимальное количество сообщений в ответе (от 1 до 100, по умолчанию 50) (опционально)
	Count *int `json:"-" query:"count,omitempty"`
}