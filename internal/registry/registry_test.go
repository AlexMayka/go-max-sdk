package registry

import (
	"strings"
	"testing"

	"github.com/AlexMayka/go-max-sdk/internal/router"
	"github.com/AlexMayka/go-max-sdk/internal/types"
)

func createTestMiddleware(name string, executionOrder *[]string) types.Middleware {
	return func(next types.Handler) types.Handler {
		return func(ctx *types.BotContext) {
			*executionOrder = append(*executionOrder, name+" start")
			next(ctx)
			*executionOrder = append(*executionOrder, name+" end")
		}
	}
}

func createTestHandler(name string, executionOrder *[]string) types.Handler {
	return func(ctx *types.BotContext) {
		*executionOrder = append(*executionOrder, name+" handler")
	}
}

func TestMiddlewareExecutionOrder(t *testing.T) {
	var executionOrder []string

	root := router.NewRouter("", nil)
	root.Use(createTestMiddleware("Root1", &executionOrder))
	root.Use(createTestMiddleware("Root2", &executionOrder))

	group := root.Group("api")
	group.Use(createTestMiddleware("Group1", &executionOrder))
	group.Use(createTestMiddleware("Group2", &executionOrder))

	route := group.OnMessage("test", createTestHandler("Test", &executionOrder))
	route.Use(createTestMiddleware("Route1", &executionOrder))
	route.Use(createTestMiddleware("Route2", &executionOrder))

	registry := BuildFrom(root)

	handlers, _ := registry.GetHandlers("", types.EventMessage)
	if len(handlers) == 0 {
		t.Fatal("No handlers found")
	}

	wrappedHandler := handlers[0].Handler

	wrappedHandler(&types.BotContext{})

	expected := []string{
		"Root1 start",
		"Root2 start",
		"Group1 start",
		"Group2 start",
		"Route1 start",
		"Route2 start",
		"Test handler",
		"Route2 end",
		"Route1 end",
		"Group2 end",
		"Group1 end",
		"Root2 end",
		"Root1 end",
	}

	if len(executionOrder) != len(expected) {
		t.Fatalf("Expected %d events, got %d", len(expected), len(executionOrder))
	}

	for i, expectedEvent := range expected {
		if executionOrder[i] != expectedEvent {
			t.Errorf("Event %d: expected '%s', got '%s'", i, expectedEvent, executionOrder[i])
		}
	}

	t.Logf("Execution order: %s", strings.Join(executionOrder, " → "))
}

func TestMultipleRoutes(t *testing.T) {
	var order1, order2 []string

	root := router.NewRouter("", nil)
	root.Use(createTestMiddleware("Root", &order1))
	root.Use(createTestMiddleware("Root", &order2))

	route1 := root.OnCommand("start", createTestHandler("Handler1", &order1))
	route1.Use(createTestMiddleware("Route1", &order1))

	route2 := root.OnMessage("help", createTestHandler("Handler2", &order2))
	route2.Use(createTestMiddleware("Route2", &order2))

	registry := BuildFrom(root)

	handlers, _ := registry.GetHandlers("", types.EventCommand)
	if len(handlers) > 0 {
		handlers[0].Handler(&types.BotContext{})
	}

	handlers, _ = registry.GetHandlers("", types.EventMessage)
	if len(handlers) > 0 {
		handlers[0].Handler(&types.BotContext{})
	}

	if len(order1) == 0 || len(order2) == 0 {
		t.Error("Routes were not executed")
	}

	t.Logf("Route 1 order: %s", strings.Join(order1, " → "))
	t.Logf("Route 2 order: %s", strings.Join(order2, " → "))
}

func TestStateHandling(t *testing.T) {
	root := router.NewRouter("", nil)

	route1 := root.OnMessage("test", func(ctx *types.BotContext) {})
	route1.UseState("auth")

	root.OnMessage("help", func(ctx *types.BotContext) {})

	registry := BuildFrom(root)

	if !registry.HasState("auth") {
		t.Error("State 'auth' not found in dispatch map")
	}

	if !registry.HasState("") {
		t.Error("Default state '' not found in dispatch map")
	}

	t.Logf("Dispatch map states: %v", registry.GetStates())
}

func BenchmarkDispatchMapBuild(b *testing.B) {
	root := router.NewRouter("", nil)

	for i := 0; i < 10; i++ {
		group := root.Group("group")
		group.Use(func(next types.Handler) types.Handler { return next })

		for j := 0; j < 10; j++ {
			route := group.OnMessage("test", func(ctx *types.BotContext) {})
			route.Use(func(next types.Handler) types.Handler { return next })
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildFrom(root)
	}
}

func BenchmarkDispatchMapAccess(b *testing.B) {
	root := router.NewRouter("", nil)
	root.OnMessage("test", func(ctx *types.BotContext) {})

	registry := BuildFrom(root)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = registry.GetHandlers("", types.EventMessage)
	}
}