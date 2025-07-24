package types

type Router interface {
	Group(prefix string) Router
	Use(middlewares ...Middleware) Router
	UseState(state string) Router

	OnCommand(cmd string, handler Handler) Route
	OnMessage(msg string, handler Handler) Route
	OnRegex(regex string, handler Handler) Route
	OnCallback(call string, handler Handler) Route
	OnStarted(handler Handler) Route
	Any(handler Handler) Route

	GetRoutes() []Route
	GetChildren() []Router
	GetMiddlewares() []Middleware
	GetState() string
	GetParent() Router
}

type Route interface {
	UseState(state string) Route
	Use(middlewares ...Middleware) Route

	GetType() RouteType
	GetPattern() *string
	GetHandler() Handler
	GetMiddleware() []Middleware
	GetState() string
}

type RouteType int

const (
	RouteMsg RouteType = iota
	RouteRegex
	RouteCommand
	RouteCallback
	RouteBotStarted
	RouteAny
)
