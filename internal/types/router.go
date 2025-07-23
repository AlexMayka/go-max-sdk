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
}

type Route interface {
	UseState(state string) Route
	Use(middlewares ...Middleware) Route
}
