// Package router provides hierarchical message routing with middleware support.
// It implements a fluent API for building route trees that get compiled into dispatch maps.
package router

import (
	"fmt"
	"github.com/AlexMayka/go-max-sdk/internal/core"
	"regexp"
)

// Router represents a node in the routing tree with support for grouping, middleware, and state.
// It provides a fluent API for building complex routing structures that inherit middleware and state.
type Router struct {
	Prefix   string
	Parent   core.Router
	Children []core.Router

	Routes      []core.Route
	Middlewares []core.Middleware
	State       string
	AnyMsg      core.Route
}

// NewRouter creates a new router with the specified prefix and parent.
// Child routers inherit middleware and state from their parents.
func NewRouter(prefix string, parent core.Router) core.Router {
	return &Router{
		Prefix:   prefix,
		Parent:   parent,
		Children: make([]core.Router, 0),

		Routes:      make([]core.Route, 0),
		Middlewares: make([]core.Middleware, 0),
		State:       "",
	}
}

// DefaultRouter creates a root router with no prefix or parent.
// This is typically used as the main router for a bot.
func DefaultRouter() core.Router {
	return NewRouter("", nil)
}

// Group creates a child router with the specified prefix.
// Child routers inherit middleware and state from their parent.
func (r *Router) Group(prefix string) core.Router {
	router := NewRouter(prefix, r)
	r.Children = append(r.Children, router)
	return router
}

// Use adds middleware to this router that will be inherited by all child routers and routes.
// Middleware are executed in the order they are added.
func (r *Router) Use(middlewares ...core.Middleware) core.Router {
	r.Middlewares = append(r.Middlewares, middlewares...)
	return r
}

// UseState sets the FSM state for this router and all its routes.
// Routes will only match messages from users in this state.
func (r *Router) UseState(state string) core.Router {
	r.State = state
	return r
}

// OnStarted registers a handler for bot startup events.
func (r *Router) OnStarted(handler core.Handler) core.Route {
	return r.addRoute(core.EventBotStarted, core.MatchBotStarted, handler, nil)
}

// OnMessage registers a handler for exact message text matches.
func (r *Router) OnMessage(msg string, handler core.Handler) core.Route {
	return r.addRoute(core.EventMessage, core.MatchExact, handler, &msg)
}

// OnRegex registers a handler for messages matching the given regular expression.
// The regex is compiled at registration time and panics on invalid patterns.
func (r *Router) OnRegex(regex string, handler core.Handler) core.Route {
	return r.addRouteRegex(core.EventMessage, core.MatchRegex, handler, regex)
}

// OnPrefix registers a handler for messages starting with the specified prefix.
func (r *Router) OnPrefix(prefix string, handler core.Handler) core.Route {
	return r.addRoute(core.EventMessage, core.MatchPrefix, handler, &prefix)
}

// OnSuffix registers a handler for messages ending with the specified suffix.
func (r *Router) OnSuffix(suffix string, handler core.Handler) core.Route {
	return r.addRoute(core.EventMessage, core.MatchSuffix, handler, &suffix)
}

// OnContains registers a handler for messages containing the specified substring.
func (r *Router) OnContains(sub string, handler core.Handler) core.Route {
	return r.addRoute(core.EventMessage, core.MatchContains, handler, &sub)
}

// OnCommand registers a handler for bot commands (messages starting with '/').
func (r *Router) OnCommand(cmd string, handler core.Handler) core.Route {
	return r.addRoute(core.EventCommand, core.MatchCommand, handler, &cmd)
}

// OnCallback registers a handler for inline keyboard callback data.
func (r *Router) OnCallback(call string, handler core.Handler) core.Route {
	return r.addRoute(core.EventCallback, core.MatchCallback, handler, &call)
}

// Any registers a catch-all handler that matches any message.
// This is typically used as a fallback when no other routes match.
func (r *Router) Any(handler core.Handler) core.Route {
	return r.addRoute(core.EventAny, core.MatchAny, handler, nil)
}

// addRoute creates and registers a new route with the specified parameters.
func (r *Router) addRoute(eventType core.EventType, matchType core.MatchType, handler core.Handler, pattern *string) core.Route {
	prefix := r.Prefix
	if pattern != nil {
		prefix = prefix + ":" + *pattern
	}

	rout := NewRoute(prefix, eventType, matchType, handler, pattern, r.State, nil)
	if eventType == core.EventAny {
		r.AnyMsg = rout
		return rout
	}

	r.Routes = append(r.Routes, rout)
	return rout
}

// addRouteRegex creates and registers a regex route with compiled pattern.
// Panics if the regex pattern is invalid.
func (r *Router) addRouteRegex(eventType core.EventType, matchType core.MatchType, handler core.Handler, pattern string) core.Route {
	prefix := r.Prefix
	if pattern != "" {
		prefix = prefix + ":" + pattern
	}

	regx, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("regex compile error: %s", err.Error()))
	}

	rout := NewRoute(prefix, eventType, matchType, handler, &pattern, r.State, regx)
	r.Routes = append(r.Routes, rout)
	return rout
}

func (r *Router) GetRoutes() []core.Route {
	routes := make([]core.Route, len(r.Routes))
	for i, route := range r.Routes {
		routes[i] = route
	}
	return routes
}

func (r *Router) GetChildren() []core.Router {
	children := make([]core.Router, len(r.Children))
	for i, child := range r.Children {
		children[i] = child
	}
	return children
}

func (r *Router) GetMiddlewares() []core.Middleware {
	return r.Middlewares
}

func (r *Router) GetState() string {
	return r.State
}

func (r *Router) GetParent() core.Router {
	if r.Parent == nil {
		return nil
	}
	return r.Parent
}

func (r *Router) GetAnyMsg() core.Route {
	return r.AnyMsg
}
