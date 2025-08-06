// Package registry compiles router trees into a flat dispatch map with middleware chains.
// It transforms the hierarchical router structure into an optimized lookup table for fast handler resolution.
package registry

import (
	"github.com/AlexMayka/go-max-sdk/internal/core"
)

// dispatchMap provides O(1) lookup for handlers by state and event type.
// Structure: state -> event_type -> []handlers_with_middleware
type dispatchMap map[string]map[core.EventType][]core.RouteHandler

// routeRegistry implements core.RouteRegistry interface.
// It stores the compiled dispatch map for fast handler lookup during message processing.
type routeRegistry struct {
	dispatch dispatchMap
}

// BuildFrom compiles a router tree into a route registry with optimized dispatch mapping.
// It recursively traverses the router tree, collects middleware chains, and builds handler lookups.
// The compilation process:
//   1. Collects middleware from parent routers (inheritance)
//   2. Wraps each handler with its middleware chain (onion pattern)
//   3. Creates flat dispatch map for O(1) handler lookup
func BuildFrom(router core.Router) core.RouteRegistry {
	r := &routeRegistry{make(dispatchMap)}
	r.buildFromRouter(router)
	return r
}

// GetHandlers retrieves all handlers for a specific state and event type.
// Returns handlers with pre-wrapped middleware chains ready for execution.
func (r *routeRegistry) GetHandlers(state string, event core.EventType) ([]core.RouteHandler, bool) {
	if _, ok := r.dispatch[state]; !ok {
		return nil, false
	}
	handlers := r.dispatch[state][event]
	return handlers, len(handlers) > 0
}

// HasState checks if the registry contains handlers for the specified state.
func (r *routeRegistry) HasState(state string) bool {
	_, exists := r.dispatch[state]
	return exists
}

// GetStates returns all states that have registered handlers.
// Useful for debugging and introspection.
func (r *routeRegistry) GetStates() []string {
	states := make([]string, 0, len(r.dispatch))
	for state := range r.dispatch {
		states = append(states, state)
	}
	return states
}

// buildFromRouter recursively processes a router tree and builds the dispatch map.
// It handles middleware inheritance, route compilation, and child router processing.
func (r *routeRegistry) buildFromRouter(router core.Router) {
	middlewares := make([]core.Middleware, 0, len(router.GetMiddlewares()))
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
		r.addValue(handler, core.EventAny, anyRoute)
	}

	for _, child := range router.GetChildren() {
		r.buildFromRouter(child)
	}
}

// addValue adds a compiled handler to the dispatch map for the specified state and event.
func (r *routeRegistry) addValue(h core.Handler, event core.EventType, route core.Route) {
	routeHandler := core.RouteHandler{Match: route.GetMatchType(), Handler: h, Route: route}
	state := route.GetState()
	if _, ok := r.dispatch[state]; !ok {
		r.dispatch[state] = make(map[core.EventType][]core.RouteHandler)
	}

	r.dispatch[state][event] = append(r.dispatch[state][event], routeHandler)
}

// collectMiddleware recursively collects middleware from parent routers.
// It traverses up the router hierarchy to ensure proper middleware inheritance.
// Middleware execution order: parent -> child -> route-specific
func (r *routeRegistry) collectMiddleware(router core.Router, middlewares []core.Middleware) []core.Middleware {
	if router.GetParent() != nil {
		middlewares = r.collectMiddleware(router.GetParent(), middlewares)
	}

	middlewares = append(middlewares, router.GetMiddlewares()...)
	return middlewares
}

// wrapHandler applies middleware chain to a handler using the onion pattern.
// Middleware are applied in reverse order to ensure correct execution sequence:
// - First middleware added wraps the outermost layer
// - Last middleware added wraps the handler directly
// Execution flows: outer -> inner -> handler -> inner -> outer
func (r *routeRegistry) wrapHandler(handler core.Handler, middlewares []core.Middleware) core.Handler {
	for index := len(middlewares) - 1; index >= 0; index-- {
		middleware := middlewares[index]
		handler = middleware(handler)
	}

	return handler
}
