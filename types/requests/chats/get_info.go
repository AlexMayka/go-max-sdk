package chats

// GetInfo представляет параметры запроса информации о чате (GET /chats/{chatId})
type GetInfo struct {
	// ID запрашиваемого чата
	ChatID int64 `json:"-" path:"chatId"`
}