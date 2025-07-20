package subscriptions

// Subscribe представляет запрос на подписку на обновления (POST /subscriptions)
type Subscribe struct {
	// URL HTTP(S)-эндпойнта вашего бота. Должен начинаться с http(s)://
	URL string `json:"url"`
	
	// Список типов обновлений, которые ваш бот хочет получать. Для полного списка типов см. объект Update (опционально)
	UpdateTypes []string `json:"update_types,omitempty"`
	
	// Cекрет, который должен быть отправлен в заголовке X-Max-Bot-Api-Secret в каждом запросе Webhook (от 5 до 256 символов) (опционально)
	Secret *string `json:"secret,omitempty"`
}