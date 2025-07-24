package types

type RouteRegistry interface {
	GetHandlers(state string, event EventType) ([]RouteHandler, bool)
	HasState(state string) bool
	GetStates() []string
}

type EventType int

const (
	EventMessage EventType = iota
	EventCommand
	EventCallback
	EventBotStarted
	EventAny
)

type RouteHandler struct {
	Handler Handler
	Route   Route
}

var TypeRouterToEvent = map[RouteType]EventType{
	RouteMsg:        EventMessage,
	RouteRegex:      EventMessage,
	RouteCommand:    EventCommand,
	RouteCallback:   EventCallback,
	RouteBotStarted: EventBotStarted,
	RouteAny:        EventAny,
}