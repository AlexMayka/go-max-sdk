package messages

// GetVideoInfo представляет запрос информации о видео (GET /videos/{videoToken})
type GetVideoInfo struct {
	// Токен видео-вложения
	VideoToken string `json:"-" path:"videoToken"`
}