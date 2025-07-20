package uploads

import "github.com/AlexMayka/go-max-sdk/types/models"

// GetUploadURL представляет URL для загрузки файла (ответ на POST /uploads)
type GetUploadURL struct {
	models.UploadInfo
}