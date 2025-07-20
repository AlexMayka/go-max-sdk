package bots

import "github.com/AlexMayka/go-max-sdk/types/models"

// UpdateInfo представляет запрос на изменение информации о боте (PATCH /me)
type UpdateInfo struct {
	// Отображаемое имя бота (от 1 до 64 символов) (опционально)
	FirstName *string `json:"first_name,omitempty"`
	
	// Отображаемое второе имя бота (от 1 до 64 символов) (опционально)
	LastName *string `json:"last_name,omitempty"`
	
	// Поле устарело, скоро будет удалено. Используйте first_name (опционально)
	Name *string `json:"name,omitempty"`
	
	// Описание бота (от 1 до 16000 символов) (опционально)
	Description *string `json:"description,omitempty"`
	
	// Команды, поддерживаемые ботом. Чтобы удалить все команды, передайте пустой список (до 32 элементов) (опционально)
	Commands []*models.BotCommand `json:"commands,omitempty"`
	
	// Запрос на установку фото бота (опционально)
	Photo *models.PhotoAttachmentRequestPayload `json:"photo,omitempty"`
}
