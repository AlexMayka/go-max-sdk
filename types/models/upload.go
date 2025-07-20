package models

// UploadType представляет тип загружаемого файла
type UploadType string

const (
	// UploadTypeImage изображение
	UploadTypeImage UploadType = "image"
	// UploadTypeVideo видео
	UploadTypeVideo UploadType = "video"
	// UploadTypeAudio аудио
	UploadTypeAudio UploadType = "audio"
	// UploadTypeFile файл
	UploadTypeFile UploadType = "file"
)

// UploadInfo представляет информацию для загрузки файла
type UploadInfo struct {
	// URL для загрузки файла
	URL string `json:"url"`
	
	// Видео- или аудио-токен для отправки сообщения (опционально)
	Token *string `json:"token,omitempty"`
}