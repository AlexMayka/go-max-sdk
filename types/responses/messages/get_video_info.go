package messages

import "github.com/AlexMayka/go-max-sdk/types/models"

// GetVideoInfo представляет информацию о видео (ответ на GET /videos/{videoToken})
type GetVideoInfo struct {
	models.VideoInfo
}