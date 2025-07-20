package chats

// SendAction представляет запрос на отправку действия в чат (POST /chats/{chatId}/actions)
type SendAction struct {
	// ID чата
	ChatID int64 `json:"-" path:"chatId"`
	
	// Действие, отправляемое участникам чата
	Action SenderAction `json:"action"`
}

// SenderAction представляет действие отправителя
type SenderAction string

const (
	// SenderActionTypingOn бот набирает сообщение
	SenderActionTypingOn SenderAction = "typing_on"
	// SenderActionSendingPhoto бот отправляет фото
	SenderActionSendingPhoto SenderAction = "sending_photo"
	// SenderActionSendingVideo бот отправляет видео
	SenderActionSendingVideo SenderAction = "sending_video"
	// SenderActionSendingAudio бот отправляет аудиофайл
	SenderActionSendingAudio SenderAction = "sending_audio"
	// SenderActionSendingFile бот отправляет файл
	SenderActionSendingFile SenderAction = "sending_file"
	// SenderActionMarkSeen бот помечает сообщения как прочитанные
	SenderActionMarkSeen SenderAction = "mark_seen"
)