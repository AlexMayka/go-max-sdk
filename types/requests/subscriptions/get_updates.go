package subscriptions

// GetUpdates представляет запрос на получение обновлений через long polling (GET /updates)
type GetUpdates struct {
	// Максимальное количество обновлений для получения (от 1 до 1000, по умолчанию 100) (опционально)
	Limit *int `json:"-" query:"limit,omitempty"`
	
	// Тайм-аут в секундах для долгого опроса (от 0 до 90, по умолчанию 30) (опционально)
	Timeout *int `json:"-" query:"timeout,omitempty"`
	
	// Если передан, бот получит обновления, которые еще не были получены. Если не передан, получит все новые обновления (опционально)
	Marker *int64 `json:"-" query:"marker,omitempty"`
	
	// Список типов обновлений, которые бот хочет получить (например, message_created, message_callback) (опционально)
	Types []string `json:"-" query:"types,omitempty"`
}