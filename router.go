package maxsdk

import (
	"github.com/AlexMayka/go-max-sdk/internal/router"
	"github.com/AlexMayka/go-max-sdk/internal/types"
)

type Router = types.Router
type Route = types.Route

func DefaultRouter() Router {
	return router.DefaultRouter()
}
