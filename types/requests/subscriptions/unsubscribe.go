package subscriptions

// Unsubscribe представляет запрос на отписку от обновлений (DELETE /subscriptions)
type Unsubscribe struct {
	// URL, который нужно удалить из подписок на WebHook
	URL string `json:"-" query:"url"`
}