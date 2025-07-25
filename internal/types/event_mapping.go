package types

import "github.com/AlexMayka/go-max-sdk/types/models"

var UpdateTypeToEvent = map[models.UpdateType]EventType{
	models.UpdateTypeMessageCreated:  EventMessage,
	models.UpdateTypeMessageCallback: EventCallback,
	models.UpdateBotStarted:          EventBotStarted,

	models.UpdateTypeMessageEdited:  EventNone,
	models.UpdateTypeMessageRemoved: EventNone,
	models.UpdateBotAdded:           EventNone,
	models.UpdateBotRemoved:         EventNone,
	models.UpdateUserAdded:          EventNone,
	models.UpdateUserRemoved:        EventNone,
	models.UpdateChatTitleChanged:   EventNone,
	models.UpdateMessageChatCreated: EventNone,
}
