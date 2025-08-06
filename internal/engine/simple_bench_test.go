package engine

import (
	"github.com/AlexMayka/go-max-sdk/internal/core"
	"testing"

	"github.com/AlexMayka/go-max-sdk/internal/registry"
	"github.com/AlexMayka/go-max-sdk/internal/router"
)

func BenchmarkBasicOperations(b *testing.B) {
	b.Run("RegistryLookup", func(b *testing.B) {
		root := router.NewRouter("", nil)
		root.OnMessage("test", func(ctx *core.BotContext) {})
		root.OnCommand("start", func(ctx *core.BotContext) {})

		reg := registry.BuildFrom(root)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_, _ = reg.GetHandlers("", core.EventMessage)
		}
	})

	b.Run("MiddlewareChain", func(b *testing.B) {
		var handler core.Handler = func(ctx *core.BotContext) {}

		for i := 0; i < 5; i++ {
			middleware := func(next core.Handler) core.Handler {
				return func(ctx *core.BotContext) {
					next(ctx)
				}
			}
			handler = middleware(handler)
		}

		ctx := &core.BotContext{}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			handler(ctx)
		}
	})

	b.Run("InterfaceCall", func(b *testing.B) {
		root := router.NewRouter("", nil)
		root.OnMessage("test", func(ctx *core.BotContext) {})

		var reg core.RouteRegistry = registry.BuildFrom(root)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_, _ = reg.GetHandlers("", core.EventMessage)
		}
	})

	b.Run("LargeRouterBuild", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			root := router.NewRouter("", nil)

			for j := 0; j < 100; j++ {
				root.OnMessage("test", func(ctx *core.BotContext) {})
			}

			_ = registry.BuildFrom(root)
		}
	})
}
