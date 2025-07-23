package router

import "github.com/AlexMayka/go-max-sdk/internal/types"

type RouteHandler struct {
	Handler types.Handler
	Route   *Route
}
