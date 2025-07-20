package messages

// Delete представляет запрос на удаление сообщения (DELETE /messages)
type Delete struct {
	// ID удаляемого сообщения (от 1 символа)
	MessageID string `json:"-" query:"message_id"`
}