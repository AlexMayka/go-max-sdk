package chats

// GetMyMembership представляет запрос информации о членстве бота в чате (GET /chats/{chatId}/members/me)
type GetMyMembership struct {
	// ID чата
	ChatID int64 `json:"-" path:"chatId"`
}