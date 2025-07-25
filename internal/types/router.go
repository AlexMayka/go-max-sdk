package types

import "regexp"

type Router interface {
	Group(prefix string) Router
	Use(middlewares ...Middleware) Router
	UseState(state string) Router

	OnCommand(cmd string, handler Handler) Route
	OnMessage(msg string, handler Handler) Route
	OnRegex(regex string, handler Handler) Route
	OnPrefix(prefix string, handler Handler) Route
	OnSuffix(suffix string, handler Handler) Route
	OnContains(sub string, handler Handler) Route
	OnCallback(call string, handler Handler) Route
	OnStarted(handler Handler) Route
	Any(handler Handler) Route

	GetRoutes() []Route
	GetChildren() []Router
	GetMiddlewares() []Middleware
	GetState() string
	GetParent() Router
	GetAnyMsg() Route
}

type Route interface {
	UseState(state string) Route
	Use(middlewares ...Middleware) Route

	GetPrefix() string
	GetMatchType() MatchType
	GetType() EventType
	GetPattern() *string
	GetHandler() Handler
	GetMiddleware() []Middleware
	GetState() string
	GetCompiledRegex() *regexp.Regexp
}
