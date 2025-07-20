package bots

import "github.com/AlexMayka/go-max-sdk/types/models"

// GetInfo представляет информацию о текущем боте (ответ на GET /me)
type GetInfo struct {
	models.BotInfo
}