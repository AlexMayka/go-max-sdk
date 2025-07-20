package messages

// GetByID представляет запрос сообщения по ID (GET /messages/{messageId})
type GetByID struct {
	// ID сообщения (mid), чтобы получить одно сообщение в чате
	MessageID string `json:"-" path:"messageId"`
}