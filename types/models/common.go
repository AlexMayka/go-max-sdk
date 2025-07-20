package models

// SuccessResponse представляет стандартный ответ об успехе операции
type SuccessResponse struct {
	// true, если запрос был успешным, false в противном случае
	Success bool `json:"success"`
	
	// Объяснительное сообщение, если результат не был успешным (опционально)
	Message *string `json:"message,omitempty"`
}