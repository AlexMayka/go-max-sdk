package subscriptions

import "github.com/AlexMayka/go-max-sdk/types/models"

// Unsubscribe представляет результат отписки от обновлений (ответ на DELETE /subscriptions)
type Unsubscribe struct {
	models.SuccessResponse
}