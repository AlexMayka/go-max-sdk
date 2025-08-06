package router

import (
	"github.com/AlexMayka/go-max-sdk/internal/core"
	"regexp"
)

// Route represents a single route with its matching criteria, handler, and middleware.
// Routes can have their own middleware that is combined with router middleware.
type Route struct {
	Prefix     string
	Type       core.EventType
	MatchType  core.MatchType
	Handler    core.Handler
	Middleware []core.Middleware

	Pattern       *string
	CompiledRegex *regexp.Regexp

	State string
}

// NewRoute creates a new route with the specified parameters.
// The route inherits state from its parent router but can override it.
func NewRoute(prefix string, eventType core.EventType, matchType core.MatchType, handler core.Handler, pattern *string, state string, regexp *regexp.Regexp) core.Route {
	return &Route{
		Prefix:        prefix,
		Type:          eventType,
		MatchType:     matchType,
		Handler:       handler,
		Middleware:    make([]core.Middleware, 0),
		Pattern:       pattern,
		State:         state,
		CompiledRegex: regexp,
	}
}

// Use adds middleware specific to this route.
// Route middleware is executed after router middleware in the chain.
func (r *Route) Use(md ...core.Middleware) core.Route {
	r.Middleware = append(r.Middleware, md...)
	return r
}

// UseState overrides the FSM state for this specific route.
func (r *Route) UseState(state string) core.Route {
	r.State = state
	return r
}

func (r *Route) GetMatchType() core.MatchType {
	return r.MatchType
}

func (r *Route) GetPrefix() string {
	return r.Prefix
}

func (r *Route) GetType() core.EventType {
	return r.Type
}

func (r *Route) GetPattern() *string {
	return r.Pattern
}

func (r *Route) GetHandler() core.Handler {
	return r.Handler
}

func (r *Route) GetMiddleware() []core.Middleware {
	return r.Middleware
}

func (r *Route) GetState() string {
	return r.State
}

func (r *Route) GetCompiledRegex() *regexp.Regexp {
	return r.CompiledRegex
}
