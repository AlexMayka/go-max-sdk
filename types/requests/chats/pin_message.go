package chats

// PinMessage представляет запрос на закрепление сообщения (PUT /chats/{chatId}/pin)
type PinMessage struct {
	// ID чата, где должно быть закреплено сообщение
	ChatID int64 `json:"-" path:"chatId"`
	
	// ID сообщения, которое нужно закрепить. Соответствует полю Message.body.mid
	MessageID string `json:"message_id"`
	
	// Если true, участники получат уведомление с системным сообщением о закреплении (по умолчанию true) (опционально)
	Notify *bool `json:"notify,omitempty"`
}