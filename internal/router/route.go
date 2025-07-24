package router

import "github.com/AlexMayka/go-max-sdk/internal/types"

type Route struct {
	Prefix     string
	Type       types.RouteType
	Handler    types.Handler
	Middleware []types.Middleware
	Pattern    *string
	State      string
}

func NewRoute(prefix string, router types.RouteType, handler types.Handler, pattern *string, state string) types.Route {
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

func (r *Route) GetType() types.RouteType {
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
