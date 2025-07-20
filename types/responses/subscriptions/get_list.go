package subscriptions

import "github.com/AlexMayka/go-max-sdk/types/models"

// GetList представляет список подписок (ответ на GET /subscriptions)
type GetList struct {
	// Список текущих подписок
	Subscriptions []models.Subscription `json:"subscriptions"`
}