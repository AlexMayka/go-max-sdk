package maxsdk

import "github.com/AlexMayka/go-max-sdk/types/models"

// KeyboardBuilder предоставляет удобный способ создания inline-клавиатур
type KeyboardBuilder struct {
	rows [][]models.InlineKeyboardButton
}

// NewKeyboard создает новый builder для клавиатуры
func NewKeyboard() *KeyboardBuilder {
	return &KeyboardBuilder{
		rows: make([][]models.InlineKeyboardButton, 0),
	}
}

// Row добавляет новый ряд кнопок
func (kb *KeyboardBuilder) Row() *RowBuilder {
	return &RowBuilder{
		keyboard: kb,
		buttons:  make([]models.InlineKeyboardButton, 0),
	}
}

// Build собирает клавиатуру
func (kb *KeyboardBuilder) Build() *models.InlineKeyboard {
	return &models.InlineKeyboard{
		Buttons: kb.rows,
	}
}

// RowBuilder предоставляет удобный способ создания ряда кнопок
type RowBuilder struct {
	keyboard *KeyboardBuilder
	buttons  []models.InlineKeyboardButton
}

// Button добавляет callback кнопку
func (rb *RowBuilder) Button(text, payload string) *RowBuilder {
	button := models.InlineKeyboardButton{
		Type:    models.InlineKeyboardButtonTypeCallback,
		Text:    text,
		Payload: &payload,
	}
	rb.buttons = append(rb.buttons, button)
	return rb
}

// ButtonWithIntent добавляет callback кнопку с намерением
func (rb *RowBuilder) ButtonWithIntent(text, payload string, intent models.ButtonIntent) *RowBuilder {
	button := models.InlineKeyboardButton{
		Type:    models.InlineKeyboardButtonTypeCallback,
		Text:    text,
		Payload: &payload,
		Intent:  &intent,
	}
	rb.buttons = append(rb.buttons, button)
	return rb
}

// PositiveButton добавляет положительную callback кнопку
func (rb *RowBuilder) PositiveButton(text, payload string) *RowBuilder {
	intent := models.ButtonIntentPositive
	return rb.ButtonWithIntent(text, payload, intent)
}

// NegativeButton добавляет негативную callback кнопку
func (rb *RowBuilder) NegativeButton(text, payload string) *RowBuilder {
	intent := models.ButtonIntentNegative
	return rb.ButtonWithIntent(text, payload, intent)
}

// URLButton добавляет кнопку-ссылку
func (rb *RowBuilder) URLButton(text, url string) *RowBuilder {
	button := models.InlineKeyboardButton{
		Type:    models.InlineKeyboardButtonTypeLink,
		Text:    text,
		Payload: &url,
	}
	rb.buttons = append(rb.buttons, button)
	return rb
}

// MessageButton добавляет кнопку, которая отправляет текстовое сообщение
func (rb *RowBuilder) MessageButton(text, message string) *RowBuilder {
	button := models.InlineKeyboardButton{
		Type:    models.InlineKeyboardButtonTypeMessage,
		Text:    text,
		Payload: &message,
	}
	rb.buttons = append(rb.buttons, button)
	return rb
}

// ContactButton добавляет кнопку запроса контакта
func (rb *RowBuilder) ContactButton(text string) *RowBuilder {
	button := models.InlineKeyboardButton{
		Type: models.InlineKeyboardButtonTypeRequestContact,
		Text: text,
	}
	rb.buttons = append(rb.buttons, button)
	return rb
}

// LocationButton добавляет кнопку запроса местоположения
func (rb *RowBuilder) LocationButton(text string) *RowBuilder {
	button := models.InlineKeyboardButton{
		Type: models.InlineKeyboardButtonTypeRequestGeoLocation,
		Text: text,
	}
	rb.buttons = append(rb.buttons, button)
	return rb
}

// AppButton добавляет кнопку открытия мини-приложения
func (rb *RowBuilder) AppButton(text, appPayload string) *RowBuilder {
	button := models.InlineKeyboardButton{
		Type:    models.InlineKeyboardButtonTypeOpenApp,
		Text:    text,
		Payload: &appPayload,
	}
	rb.buttons = append(rb.buttons, button)
	return rb
}

// End завершает создание ряда и добавляет его в клавиатуру
func (rb *RowBuilder) End() *KeyboardBuilder {
	rb.keyboard.rows = append(rb.keyboard.rows, rb.buttons)
	return rb.keyboard
}