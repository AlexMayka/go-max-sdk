package router

import (
	"fmt"
	"github.com/AlexMayka/go-max-sdk/internal/types"
	"regexp"
)

type Router struct {
	Prefix   string
	Parent   types.Router
	Children []types.Router

	Routes      []types.Route
	Middlewares []types.Middleware
	State       string
}

func NewRouter(prefix string, parent *Router) types.Router {
	return &Router{
		Prefix:   prefix,
		Parent:   parent,
		Children: make([]types.Router, 0),

		Routes:      make([]types.Route, 0),
		Middlewares: make([]types.Middleware, 0),
		State:       "",
	}
}

func DefaultRouter() types.Router {
	return NewRouter("", nil)
}

func (r *Router) Group(prefix string) types.Router {
	router := NewRouter(prefix, r)
	r.Children = append(r.Children, router)
	return router
}

func (r *Router) Use(middlewares ...types.Middleware) types.Router {
	r.Middlewares = append(r.Middlewares, middlewares...)
	return r
}

func (r *Router) UseState(state string) types.Router {
	r.State = state
	return r
}

func (r *Router) OnStarted(handler types.Handler) types.Route {
	return r.addRoute(RouteBotStarted, handler, nil)
}

func (r *Router) OnMessage(msg string, handler types.Handler) types.Route {
	return r.addRoute(RouteMsg, handler, &msg)
}

func (r *Router) OnCommand(cmd string, handler types.Handler) types.Route {
	return r.addRoute(RouteCommand, handler, &cmd)
}

func (r *Router) OnCallback(call string, handler types.Handler) types.Route {
	return r.addRoute(RouteCallback, handler, &call)
}

func (r *Router) OnRegex(regex string, handler types.Handler) types.Route {
	if _, err := regexp.Compile(regex); err != nil {
		panic(fmt.Sprintf("regex compile error: %s", err.Error()))
	}

	return r.addRoute(RouteRegex, handler, &regex)
}

func (r *Router) Any(handler types.Handler) types.Route {
	return r.addRoute(RouteAny, handler, nil)
}

func (r *Router) addRoute(router types.RouteType, handler types.Handler, pattern *string) types.Route {
	prefix := r.Prefix
	if pattern != nil {
		prefix = prefix + ":" + *pattern
	}

	route := NewRoute(prefix, router, handler, pattern, r.State)
	r.Routes = append(r.Routes, route)
	return route
}

func (r *Router) GetRoutes() []types.Route {
	routes := make([]types.Route, len(r.Routes))
	for i, route := range r.Routes {
		routes[i] = route
	}
	return routes
}

func (r *Router) GetChildren() []types.Router {
	children := make([]types.Router, len(r.Children))
	for i, child := range r.Children {
		children[i] = child
	}
	return children
}

func (r *Router) GetMiddlewares() []types.Middleware {
	return r.Middlewares
}

func (r *Router) GetState() string {
	return r.State
}

func (r *Router) GetParent() types.Router {
	if r.Parent == nil {
		return nil
	}
	return r.Parent
}
