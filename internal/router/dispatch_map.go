package router

import (
	"github.com/AlexMayka/go-max-sdk/internal/types"
)

type DispatchMap map[string]map[EventType][]RouteHandler

func NewDispatchMap(router *Router) *DispatchMap {
	d := make(DispatchMap)
	d.buildDispatchMap(router)
	return &d
}

func (d *DispatchMap) buildDispatchMap(router *Router) {
	middlewares := make([]types.Middleware, 0, len(router.Middlewares))
	middlewares = collectMiddleware(router, middlewares)

	for _, route := range router.Routes {
		routeMiddlewares := append(middlewares, route.Middleware...)
		handler := wrapHandler(route.Handler, routeMiddlewares)
		event := selectEvent(route)
		d.addValue(handler, event, route)
	}

	for _, child := range router.Children {
		d.buildDispatchMap(child)
	}
}

func (d *DispatchMap) addValue(h types.Handler, event EventType, r *Route) {
	routeHandler := RouteHandler{Handler: h, Route: r}
	if _, ok := (*d)[r.State]; !ok {
		(*d)[r.State] = make(map[EventType][]RouteHandler)
	}

	(*d)[r.State][event] = append((*d)[r.State][event], routeHandler)
}

func collectMiddleware(router *Router, middlewares []types.Middleware) []types.Middleware {
	if router.Parent != nil {
		middlewares = collectMiddleware(router.Parent, middlewares)
	}

	middlewares = append(middlewares, router.Middlewares...)

	return middlewares
}

func wrapHandler(handler types.Handler, middlewares []types.Middleware) types.Handler {
	for index := len(middlewares) - 1; index >= 0; index-- {
		middleware := middlewares[index]
		handler = middleware(handler)
	}

	return handler
}

func selectEvent(r *Route) EventType {
	return TypeRouterToEvent[r.Type]
}
