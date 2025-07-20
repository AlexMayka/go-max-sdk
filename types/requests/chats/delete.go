package chats

// Delete представляет запрос на удаление чата (DELETE /chats/{chatId})
type Delete struct {
	// ID чата
	ChatID int64 `json:"-" path:"chatId"`
}