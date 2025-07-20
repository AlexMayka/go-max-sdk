package chats

// UnpinMessage представляет запрос на удаление закрепленного сообщения (DELETE /chats/{chatId}/pin)
type UnpinMessage struct {
	// ID чата, из которого нужно удалить закрепленное сообщение
	ChatID int64 `json:"-" path:"chatId"`
}