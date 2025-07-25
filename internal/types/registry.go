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
	EventNone
)

type RouteHandler struct {
	Match   MatchType
	Handler Handler
	Route   Route
}
