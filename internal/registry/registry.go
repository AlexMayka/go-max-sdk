package registry

import (
	"github.com/AlexMayka/go-max-sdk/internal/types"
)

type dispatchMap map[string]map[types.EventType][]types.RouteHandler

type routeRegistry struct {
	dispatch dispatchMap
}

func BuildFrom(router types.Router) types.RouteRegistry {
	r := &routeRegistry{make(dispatchMap)}
	r.buildFromRouter(router)
	return r
}

func (r *routeRegistry) GetHandlers(state string, event types.EventType) ([]types.RouteHandler, bool) {
	if _, ok := r.dispatch[state]; !ok {
		return nil, false
	}
	handlers := r.dispatch[state][event]
	return handlers, len(handlers) > 0
}

func (r *routeRegistry) HasState(state string) bool {
	_, exists := r.dispatch[state]
	return exists
}

func (r *routeRegistry) GetStates() []string {
	states := make([]string, 0, len(r.dispatch))
	for state := range r.dispatch {
		states = append(states, state)
	}
	return states
}

func (r *routeRegistry) buildFromRouter(router types.Router) {
	middlewares := make([]types.Middleware, 0, len(router.GetMiddlewares()))
	middlewares = r.collectMiddleware(router, middlewares)

	for _, route := range router.GetRoutes() {
		routeMiddlewares := append(middlewares, route.GetMiddleware()...)
		handler := r.wrapHandler(route.GetHandler(), routeMiddlewares)
		eventType := route.GetType()
		r.addValue(handler, eventType, route)
	}

	if anyRoute := router.GetAnyMsg(); anyRoute != nil {
		routeMiddlewares := append(middlewares, anyRoute.GetMiddleware()...)
		handler := r.wrapHandler(anyRoute.GetHandler(), routeMiddlewares)
		r.addValue(handler, types.EventAny, anyRoute)
	}

	for _, child := range router.GetChildren() {
		r.buildFromRouter(child)
	}
}

func (r *routeRegistry) addValue(h types.Handler, event types.EventType, route types.Route) {
	routeHandler := types.RouteHandler{Match: route.GetMatchType(), Handler: h, Route: route}
	state := route.GetState()
	if _, ok := r.dispatch[state]; !ok {
		r.dispatch[state] = make(map[types.EventType][]types.RouteHandler)
	}

	r.dispatch[state][event] = append(r.dispatch[state][event], routeHandler)
}

func (r *routeRegistry) collectMiddleware(router types.Router, middlewares []types.Middleware) []types.Middleware {
	if router.GetParent() != nil {
		middlewares = r.collectMiddleware(router.GetParent(), middlewares)
	}

	middlewares = append(middlewares, router.GetMiddlewares()...)
	return middlewares
}

func (r *routeRegistry) wrapHandler(handler types.Handler, middlewares []types.Middleware) types.Handler {
	for index := len(middlewares) - 1; index >= 0; index-- {
		middleware := middlewares[index]
		handler = middleware(handler)
	}

	return handler
}
