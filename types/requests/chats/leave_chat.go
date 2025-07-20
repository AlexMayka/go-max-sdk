package chats

// LeaveChat представляет запрос на удаление бота из чата (DELETE /chats/{chatId}/members/me)
type LeaveChat struct {
	// ID чата
	ChatID int64 `json:"-" path:"chatId"`
}