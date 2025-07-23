package router

type RouteType int

const (
	RouteBotStarted RouteType = iota
	RouteMsg
	RouteCommand
	RouteCallback
	RouteRegex
	RouteAny
)

type EventType int

const (
	EventMessage EventType = iota
	EventCommand
	EventCallback
	EventBotStarted
	EventAny
)

var TypeRouterToEvent = map[RouteType]EventType{
	RouteMsg:        EventMessage,
	RouteRegex:      EventMessage,
	RouteCommand:    EventCommand,
	RouteCallback:   EventCallback,
	RouteBotStarted: EventBotStarted,
	RouteAny:        EventAny,
}
