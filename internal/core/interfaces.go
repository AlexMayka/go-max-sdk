// Package core defines the fundamental interfaces and contracts for the bot SDK.
// It provides clean abstractions for all major components without implementation details.
package core

import (
	"context"
	"github.com/AlexMayka/go-max-sdk/types/models"
	"regexp"
)

// APIClient provides the interface for making API calls to the Max messaging platform
type APIClient interface {
	Call(ctx context.Context, endpoint Endpoint, req interface{}) (interface{}, error)
}

// BotEngine interface for bot processing engine
type BotEngine interface {
	Start() error
	Stop() error
}

// Logger interface for structured logging
type Logger interface {
	Debug(component, event, details string)
	Info(component, event, details string)
	Warn(component, event, details string)
	Error(component, event, details string)
}

// Transport interface for communication with the MAX API
type Transport interface {
	Start(ctx context.Context) (<-chan *models.Update, error)
	Stop() error
}

// FSM interface for finite state machine
type FSM interface {
	CreateUser(id int64, state string) bool
	GetStateUser(id int64) (string, bool)
	SetStateUser(id int64, state string) bool
	SetValue(id int64, param string, data string) bool
	GetValue(id int64, param string) (string, bool)
}

// Router interface for message routing
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

// Route interface for individual routes
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

// RouteRegistry interface for managing routes
type RouteRegistry interface {
	GetHandlers(state string, event EventType) ([]RouteHandler, bool)
	HasState(state string) bool
	GetStates() []string
}
