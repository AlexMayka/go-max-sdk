package models

// Subscription представляет подписку на обновления через WebHook
type Subscription struct {
	// URL вебхука
	URL string `json:"url"`
	
	// Unix-время, когда была создана подписка
	Time int64 `json:"time"`
	
	// Типы обновлений, на которые подписан бот
	UpdateTypes []string `json:"update_types,omitempty"`
}