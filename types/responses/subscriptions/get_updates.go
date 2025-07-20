package subscriptions

import "github.com/AlexMayka/go-max-sdk/types/models"

// GetUpdates представляет страницу обновлений (ответ на GET /updates)
type GetUpdates struct {
	// Страница обновлений
	Updates []models.Update `json:"updates"`
	
	// Указатель на следующую страницу данных
	Marker *int64 `json:"marker,omitempty"`
}