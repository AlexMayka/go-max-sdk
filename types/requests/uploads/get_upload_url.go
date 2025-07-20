package uploads

import "github.com/AlexMayka/go-max-sdk/types/models"

// GetUploadURL представляет запрос на получение URL для загрузки файла (POST /uploads)
type GetUploadURL struct {
	// Тип загружаемого файла (enum: "image", "video", "audio", "file")
	Type models.UploadType `json:"-" query:"type"`
}