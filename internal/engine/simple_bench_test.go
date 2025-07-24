package engine

import (
	"testing"

	"github.com/AlexMayka/go-max-sdk/internal/registry"
	"github.com/AlexMayka/go-max-sdk/internal/router"
	"github.com/AlexMayka/go-max-sdk/internal/types"
)

func BenchmarkBasicOperations(b *testing.B) {
	b.Run("RegistryLookup", func(b *testing.B) {
		root := router.NewRouter("", nil)
		root.OnMessage("test", func(ctx *types.BotContext) {})
		root.OnCommand("start", func(ctx *types.BotContext) {})

		reg := registry.BuildFrom(root)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_, _ = reg.GetHandlers("", types.EventMessage)
		}
	})

	b.Run("MiddlewareChain", func(b *testing.B) {
		var handler types.Handler = func(ctx *types.BotContext) {}

		for i := 0; i < 5; i++ {
			middleware := func(next types.Handler) types.Handler {
				return func(ctx *types.BotContext) {
					next(ctx)
				}
			}
			handler = middleware(handler)
		}

		ctx := &types.BotContext{}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			handler(ctx)
		}
	})

	b.Run("InterfaceCall", func(b *testing.B) {
		root := router.NewRouter("", nil)
		root.OnMessage("test", func(ctx *types.BotContext) {})

		var reg types.RouteRegistry = registry.BuildFrom(root)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_, _ = reg.GetHandlers("", types.EventMessage)
		}
	})

	b.Run("LargeRouterBuild", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			root := router.NewRouter("", nil)

			for j := 0; j < 100; j++ {
				root.OnMessage("test", func(ctx *types.BotContext) {})
			}

			_ = registry.BuildFrom(root)
		}
	})
}
