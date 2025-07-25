package types

import (
	"context"
	"github.com/AlexMayka/go-max-sdk/types/models"
)

type BotContext struct {
	Ctx    context.Context
	Cancel context.CancelFunc

	// Данные пользователя и чата
	UserID int64
	ChatID int64
	
	// Исходное обновление
	Update *models.Update
	
	// Доступ к компонентам
	Client APIClient
	FSM    FSM
	
	// Хелперы
	Text     string  // Текст сообщения (если есть)
	Payload  *string // Payload для callback (если есть)
}
