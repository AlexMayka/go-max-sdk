package chats

// GetPinnedMessage представляет запрос на получение закрепленного сообщения (GET /chats/{chatId}/pin)
type GetPinnedMessage struct {
	// ID чата
	ChatID int64 `json:"-" path:"chatId"`
}