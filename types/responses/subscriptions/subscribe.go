package subscriptions

import "github.com/AlexMayka/go-max-sdk/types/models"

// Subscribe представляет результат подписки на обновления (ответ на POST /subscriptions)
type Subscribe struct {
	models.SuccessResponse
}