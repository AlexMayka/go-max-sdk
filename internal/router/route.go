package router

import "github.com/AlexMayka/go-max-sdk/internal/types"

type Route struct {
	Prefix     string
	Type       RouteType
	Handler    types.Handler
	Middleware []types.Middleware
	Pattern    *string
	State      string
}

func NewRoute(prefix string, router RouteType, handler types.Handler, pattern *string, state string) types.Route {
	return &Route{
		Prefix:     prefix,
		Type:       router,
		Handler:    handler,
		Middleware: make([]types.Middleware, 0),
		Pattern:    pattern,
		State:      state,
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
