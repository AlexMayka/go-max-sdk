package chats

// AddMembers представляет запрос на добавление участников в чат (POST /chats/{chatId}/members)
type AddMembers struct {
	// ID чата
	ChatID int64 `json:"-" path:"chatId"`
	
	// Массив ID пользователей для добавления в чат
	UserIDs []int64 `json:"user_ids"`
}