package models

// Message представляет сообщение в чате
type Message struct {
	// Пользователь, отправивший сообщение (опционально)
	Sender *User `json:"sender,omitempty"`
	// Получатель сообщения. Может быть пользователем или чатом
	Recipient *Recipient `json:"recipient"`
	// Время создания сообщения в формате Unix-time
	Timestamp int64 `json:"timestamp"`
	// Пересланное или ответное сообщение (опционально)
	Link *LinkedMessage `json:"link,omitempty"`
	// Содержимое сообщения. Текст + вложения. Может быть null, если сообщение содержит только пересланное сообщение
	Body *MessageBody `json:"body"`
	// Статистика сообщения (опционально)
	Stat *MessageStat `json:"stat,omitempty"`
	// Публичная ссылка на сообщение. Может быть null для диалогов или не публичных чатов (опционально)
	Url string `json:"url,omitempty"`
}

// Recipient получатель сообщения. Может быть пользователем или чатом
type Recipient struct {
	// ID чата (может быть null)
	ChatID int64 `json:"chat_id"`
	// Тип чата (enum: "chat")
	ChatType string `json:"chat_type"`
	// ID пользователя, если сообщение было отправлено пользователю (может быть null)
	UserID int64 `json:"user_id"`
}

// LinkedMessage пересланное или ответное сообщение
type LinkedMessage struct {
	// Тип связанного сообщения (enum: "forward", "reply")
	Type string `json:"type"`
	// Пользователь, отправивший сообщение (опционально)
	Sender *User `json:"sender,omitempty"`
	// Чат, в котором сообщение было изначально опубликовано. Только для пересланных сообщений (опционально)
	ChatID int64 `json:"chat_id,omitempty"`
	// Схема, представляющая тело сообщения
	Message *Message `json:"message,omitempty"`
}

// MessageBody содержимое сообщения. Текст + вложения
type MessageBody struct {
	// Уникальный ID сообщения
	Mid string `json:"mid"`
	// ID последовательности сообщения в чате
	Seq int64 `json:"seq"`
	// Текст сообщения (может быть null)
	Text string `json:"text"`
	// Вложения сообщения. Могут быть одним из типов Attachment
	Attachments []*Attachment `json:"attachments"`
	// Разметка текста сообщения (опционально)
	Markup []*MarkupElement `json:"markup,omitempty"`
}

// Attachment вложение сообщения
type Attachment struct {
	// Тип вложения
	Type string `json:"type"`
	// Полезная нагрузка вложения (может быть разных типов)
	Payload interface{} `json:"payload"`
}

// PhotoAttachmentPayload полезная нагрузка для фото-вложения
type PhotoAttachmentPayload struct {
	// ID фотографии
	PhotoID int64 `json:"photo_id"`
	// Токен для доступа к фото
	Token string `json:"token"`
	// URL фотографии
	Url string `json:"url"`
}

// MarkupElement элемент разметки текста
type MarkupElement struct {
	// Тип разметки
	Type string `json:"type"`
	// Начальная позиция в тексте
	From int32 `json:"from"`
	// Длина размеченного текста
	Length int32 `json:"length"`
}

// MessageStat статистика сообщения
type MessageStat struct {
	// Количество просмотров
	Views int64 `json:"views"`
}

// NewMessageBody представляет новое тело сообщения для отправки
type NewMessageBody struct {
	// Новый текст сообщения (до 4000 символов) (может быть null)
	Text *string `json:"text,omitempty"`

	// Вложения сообщения. Если пусто, все вложения будут удалены
	Attachments []AttachmentRequest `json:"attachments,omitempty"`

	// Ссылка на сообщение
	Link *NewMessageLink `json:"link,omitempty"`

	// Если false, участники чата не будут уведомлены (по умолчанию true) (опционально)
	Notify *bool `json:"notify,omitempty"`

	// Если установлен, текст сообщения будет форматрован данным способом (опционально)
	Format *TextFormat `json:"format,omitempty"`
}

// AttachmentRequest представляет запрос на вложение
type AttachmentRequest struct {
	// Тип вложения
	Type string `json:"type"`
	// Полезная нагрузка для запроса на прикрепление (может быть разных типов)
	Payload interface{} `json:"payload"`
}

// PhotoAttachmentRequestPayload запрос на прикрепление изображения (все поля являются взаимоисключающими)
type PhotoAttachmentRequestPayload struct {
	// URL изображения
	URL *string `json:"url,omitempty"`
	// Токен изображения
	Token *string `json:"token,omitempty"`
	// Фотографии
	Photos *PhotosPayload `json:"photos,omitempty"`
}

// InlineKeyboardAttachmentPayload полезная нагрузка для inline-клавиатуры
type InlineKeyboardAttachmentPayload struct {
	InlineKeyboard
}

// PhotosPayload представляет полезную нагрузку для фотографий
type PhotosPayload struct {
	// Токен фотографии
	Token string `json:"token"`
}

// NewMessageLink представляет ссылку на сообщение
type NewMessageLink struct {
	// Тип ссылки сообщения (enum: "forward", "reply")
	Type MessageLinkType `json:"type"`
	// ID сообщения исходного сообщения
	MID string `json:"mid"`
}

// TextFormat представляет формат текста
type TextFormat string

const (
	// TextFormatMarkdown форматирование markdown
	TextFormatMarkdown TextFormat = "markdown"
	// TextFormatHTML форматирование HTML
	TextFormatHTML TextFormat = "html"
)

// MessageLinkType представляет тип связанного сообщения
type MessageLinkType string

const (
	// MessageLinkTypeForward пересланное сообщение
	MessageLinkTypeForward MessageLinkType = "forward"
	// MessageLinkTypeReply ответ на сообщение
	MessageLinkTypeReply MessageLinkType = "reply"
)
