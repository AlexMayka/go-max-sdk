package chats

// GetAdmins представляет запрос списка администраторов чата (GET /chats/{chatId}/members/admins)
type GetAdmins struct {
	// ID чата
	ChatID int64 `json:"-" path:"chatId"`
}