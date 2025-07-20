package models

// ChatMember представляет участника чата
type ChatMember struct {
	// ID пользователя
	UserID int64 `json:"user_id"`
	
	// Отображаемое имя пользователя
	FirstName string `json:"first_name"`
	
	// Отображаемая фамилия пользователя (может быть null)
	LastName *string `json:"last_name,omitempty"`
	
	// Устаревшее поле, скоро будет удалено (опционально)
	Name *string `json:"name,omitempty"`
	
	// Уникальное публичное имя пользователя. Может быть null, если пользователь недоступен или имя не задано
	Username *string `json:"username,omitempty"`
	
	// true, если пользователь является ботом
	IsBot bool `json:"is_bot"`
	
	// Время последней активности пользователя в MAX (Unix-время в миллисекундах). Может быть неактуальным, если пользователь отключил статус "онлайн" в настройках
	LastActivityTime int64 `json:"last_activity_time"`
	
	// Описание пользователя. Может быть null, если пользователь его не заполнил (до 16000 символов) (опционально)
	Description *string `json:"description,omitempty"`
	
	// URL аватара (опционально)
	AvatarURL *string `json:"avatar_url,omitempty"`
	
	// URL аватара большего размера (опционально)
	FullAvatarURL *string `json:"full_avatar_url,omitempty"`
	
	// Время последней активности пользователя в чате. Может быть устаревшим для суперчатов (равно времени вступления)
	LastAccessTime int64 `json:"last_access_time"`
	
	// Является ли пользователь владельцем чата
	IsOwner bool `json:"is_owner"`
	
	// Является ли пользователь администратором чата
	IsAdmin bool `json:"is_admin"`
	
	// Дата присоединения к чату в формате Unix time
	JoinTime int64 `json:"join_time"`
	
	// Перечень прав пользователя
	Permissions []ChatAdminPermission `json:"permissions,omitempty"`
	
	// Заголовок, который будет показан на клиенте. Если пользователь администратор или владелец и ему не установлено это название, то поле не передается (опционально)
	Alias *string `json:"alias,omitempty"`
}

// ChatAdminPermission представляет права администратора чата
type ChatAdminPermission string

const (
	// ChatAdminPermissionReadAllMessages читать все сообщения
	ChatAdminPermissionReadAllMessages ChatAdminPermission = "read_all_messages"
	// ChatAdminPermissionAddRemoveMembers добавлять/удалять участников
	ChatAdminPermissionAddRemoveMembers ChatAdminPermission = "add_remove_members"
	// ChatAdminPermissionAddAdmins добавлять администраторов
	ChatAdminPermissionAddAdmins ChatAdminPermission = "add_admins"
	// ChatAdminPermissionChangeChatInfo изменять информацию о чате
	ChatAdminPermissionChangeChatInfo ChatAdminPermission = "change_chat_info"
	// ChatAdminPermissionPinMessage закреплять сообщения
	ChatAdminPermissionPinMessage ChatAdminPermission = "pin_message"
	// ChatAdminPermissionWrite писать сообщения
	ChatAdminPermissionWrite ChatAdminPermission = "write"
)

// ChatAdmin представляет администратора чата для назначения
type ChatAdmin struct {
	// ID пользователя
	UserID int64 `json:"user_id"`
	
	// Перечень прав администратора
	Permissions []ChatAdminPermission `json:"permissions"`
	
	// Заголовок администратора (опционально)
	Alias *string `json:"alias,omitempty"`
}