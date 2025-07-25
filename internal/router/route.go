package router

import (
	"github.com/AlexMayka/go-max-sdk/internal/types"
	"regexp"
)

type Route struct {
	Prefix     string
	Type       types.EventType
	MatchType  types.MatchType
	Handler    types.Handler
	Middleware []types.Middleware

	Pattern       *string
	CompiledRegex *regexp.Regexp

	State string
}

func NewRoute(prefix string, eventType types.EventType, matchType types.MatchType, handler types.Handler, pattern *string, state string, regexp *regexp.Regexp) types.Route {
	return &Route{
		Prefix:        prefix,
		Type:          eventType,
		MatchType:     matchType,
		Handler:       handler,
		Middleware:    make([]types.Middleware, 0),
		Pattern:       pattern,
		State:         state,
		CompiledRegex: regexp,
	}
}

func (r *Route) Use(md ...types.Middleware) types.Route {
	r.Middleware = append(r.Middleware, md...)
	return r
}

func (r *Route) UseState(state string) types.Route {
	r.State = state
	return r
}

func (r *Route) GetMatchType() types.MatchType {
	return r.MatchType
}

func (r *Route) GetPrefix() string {
	return r.Prefix
}

func (r *Route) GetType() types.EventType {
	return r.Type
}

func (r *Route) GetPattern() *string {
	return r.Pattern
}

func (r *Route) GetHandler() types.Handler {
	return r.Handler
}

func (r *Route) GetMiddleware() []types.Middleware {
	return r.Middleware
}

func (r *Route) GetState() string {
	return r.State
}

func (r *Route) GetCompiledRegex() *regexp.Regexp {
	return r.CompiledRegex
}
