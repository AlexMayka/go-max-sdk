package chats

// GetByLink представляет параметры запроса чата по ссылке (GET /chats/{chatLink})
type GetByLink struct {
	// Публичная ссылка на чат или username пользователя. Должна соответствовать регулярному выражению @?[a-zA-Z]+[a-zA-Z0-9-_]*
	ChatLink string `json:"-" path:"chatLink"`
}