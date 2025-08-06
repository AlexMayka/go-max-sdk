package core

import (
	"github.com/AlexMayka/go-max-sdk/types/models"
)

// Endpoint represents an API endpoint constant
type Endpoint string

// API endpoint constants
const (
	SendMsg            Endpoint = "SendMsg"
	EditMsg            Endpoint = "EditMsg"
	DeleteMsg          Endpoint = "DeleteMsg"
	GetSubscribeUpdate Endpoint = "GetSubscribeUpdate"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// Handler is a function that processes incoming updates
type Handler func(*BotContext)

// Middleware is a function that wraps handlers to provide cross-cutting functionality
type Middleware func(Handler) Handler

// EventType represents the type of event
type EventType int

const (
	EventMessage EventType = iota
	EventCommand
	EventCallback
	EventBotStarted
	EventAny
	EventNone
)

// MatchType represents how routes are matched
type MatchType int

const (
	MatchExact MatchType = iota
	MatchPrefix
	MatchSuffix
	MatchContains
	MatchRegex
	MatchCommand
	MatchCallback
	MatchBotStarted
	MatchAny
)

// RouteHandler represents a handler with its matching criteria
type RouteHandler struct {
	Match   MatchType
	Handler Handler
	Route   Route
}

// UpdateTypeToEvent maps update types to event types
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
