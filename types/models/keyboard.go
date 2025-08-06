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

	// Payload для кнопок (до 1024 символов)
	Payload *string `json:"payload,omitempty"`

	// Намерение кнопки
	Intent *ButtonIntent `json:"intent,omitempty"`
}

// InlineKeyboardButtonType представляет тип кнопки в inline-клавиатуре
type InlineKeyboardButtonType string

const (
	// InlineKeyboardButtonTypeCallback Callback кнопка
	InlineKeyboardButtonTypeCallback InlineKeyboardButtonType = "callback"

	// InlineKeyboardButtonTypeLink Link кнопка
	InlineKeyboardButtonTypeLink InlineKeyboardButtonType = "link"

	// InlineKeyboardButtonTypeRequestGeoLocation Request geo location кнопка
	InlineKeyboardButtonTypeRequestGeoLocation InlineKeyboardButtonType = "request_geo_location"

	// InlineKeyboardButtonTypeRequestContact Request contact кнопка
	InlineKeyboardButtonTypeRequestContact InlineKeyboardButtonType = "request_contact"

	// InlineKeyboardButtonTypeOpenApp Open app кнопка
	InlineKeyboardButtonTypeOpenApp InlineKeyboardButtonType = "open_app"

	// InlineKeyboardButtonTypeMessage Message кнопка
	InlineKeyboardButtonTypeMessage InlineKeyboardButtonType = "message"
)

// ButtonIntent представляет намерение кнопки
type ButtonIntent string

const (
	// ButtonIntentPositive положительное намерение
	ButtonIntentPositive ButtonIntent = "positive"

	// ButtonIntentNegative негативное намерение
	ButtonIntentNegative ButtonIntent = "negative"

	// ButtonIntentDefault обычное намерение
	ButtonIntentDefault ButtonIntent = "default"
)
