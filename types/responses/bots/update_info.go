package bots

import "github.com/AlexMayka/go-max-sdk/types/models"

// UpdateInfo представляет обновленную информацию о боте (ответ на PATCH /me)
type UpdateInfo struct {
	models.BotInfo
}
