package models

// InlineKeyboard представляет inline-клавиатуру
type InlineKeyboard struct {
	// Массив кнопок, сгруппированных по рядам (до 30 рядов, до 7 кнопок в ряду)
	Buttons [][]InlineKeyboardButton `json:"buttons"`
}

// InlineKeyboardButton представляет кнопку в inline-клавиатуре
type InlineKeyboardButton struct {
	// Тип кнопки
	Type InlineKeyboardButtonType `json:"type"`

	// Текст на кнопке
	Text string `json:"text"`

	// Payload для callback кнопок
	Payload *string `json:"payload,omitempty"`

	// URL для link кнопок
	URL *string `json:"url,omitempty"`
}

// InlineKeyboardButtonType представляет тип кнопки в inline-клавиатуре
type InlineKeyboardButtonType string

const (
	// InlineKeyboardButtonTypeCallback Callback кнопка - сервер MAX отправляет событие с типом message_callback
	InlineKeyboardButtonTypeCallback InlineKeyboardButtonType = "callback"

	// InlineKeyboardButtonTypeLink Link кнопка - позволяет открыть ссылку в новой вкладке
	InlineKeyboardButtonTypeLink InlineKeyboardButtonType = "link"

	// InlineKeyboardButtonTypeRequestContact Request contact кнопка - запрашивает у пользователя разрешение на доступ к контактам
	InlineKeyboardButtonTypeRequestContact InlineKeyboardButtonType = "request_contact"

	// InlineKeyboardButtonTypeRequestGeoLocation Request geo location кнопка - запрашивает у пользователя его местоположение
	InlineKeyboardButtonTypeRequestGeoLocation InlineKeyboardButtonType = "request_geo_location"

	// InlineKeyboardButtonTypeOpenApp Open app кнопка - открывает мини-приложение
	InlineKeyboardButtonTypeOpenApp InlineKeyboardButtonType = "open_app"

	// InlineKeyboardButtonTypeMessage Message кнопка - отправляет боту текстовое сообщение
	InlineKeyboardButtonTypeMessage InlineKeyboardButtonType = "message"
)
