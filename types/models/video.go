package models

// VideoInfo представляет информацию о видео
type VideoInfo struct {
	// Токен видео-вложения
	Token string `json:"token"`
	
	// URL-ы для скачивания или воспроизведения видео. Может быть null, если видео недоступно (опционально)
	URLs *VideoUrls `json:"urls,omitempty"`
	
	// Миниатюра видео (опционально)
	Thumbnail *PhotoAttachmentPayload `json:"thumbnail,omitempty"`
	
	// Ширина видео
	Width int `json:"width"`
	
	// Высота видео
	Height int `json:"height"`
	
	// Длина видео в секундах
	Duration int `json:"duration"`
}

// VideoUrls представляет URL-ы для видео
type VideoUrls struct {
	// URL видео в разрешении 1080p, если доступно
	Mp4_1080 *string `json:"mp4_1080,omitempty"`
	
	// URL видео в разрешении 720p, если доступно
	Mp4_720 *string `json:"mp4_720,omitempty"`
	
	// URL видео в разрешении 480p, если доступно
	Mp4_480 *string `json:"mp4_480,omitempty"`
	
	// URL видео в разрешении 360p, если доступно
	Mp4_360 *string `json:"mp4_360,omitempty"`
	
	// URL видео в разрешении 240p, если доступно
	Mp4_240 *string `json:"mp4_240,omitempty"`
	
	// URL видео в разрешении 144p, если доступно
	Mp4_144 *string `json:"mp4_144,omitempty"`
	
	// URL трансляции, если доступна
	HLS *string `json:"hls,omitempty"`
}