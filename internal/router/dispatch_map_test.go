package router

import (
	"strings"
	"testing"

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

	root := NewRouter("", nil).(*Router)
	root.Use(createTestMiddleware("Root1", &executionOrder))
	root.Use(createTestMiddleware("Root2", &executionOrder))

	group := root.Group("api").(*Router)
	group.Use(createTestMiddleware("Group1", &executionOrder))
	group.Use(createTestMiddleware("Group2", &executionOrder))

	route := group.OnMessage("test", createTestHandler("Test", &executionOrder)).(*Route)
	route.Use(createTestMiddleware("Route1", &executionOrder))
	route.Use(createTestMiddleware("Route2", &executionOrder))

	dispatchMap := NewDispatchMap(root)

	handlers := (*dispatchMap)[""][EventMessage]
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

	root := NewRouter("", nil).(*Router)
	root.Use(createTestMiddleware("Root", &order1))
	root.Use(createTestMiddleware("Root", &order2))

	route1 := root.OnCommand("start", createTestHandler("Handler1", &order1)).(*Route)
	route1.Use(createTestMiddleware("Route1", &order1))

	route2 := root.OnMessage("help", createTestHandler("Handler2", &order2)).(*Route)
	route2.Use(createTestMiddleware("Route2", &order2))

	dispatchMap := NewDispatchMap(root)

	handlers := (*dispatchMap)[""][EventCommand]
	if len(handlers) > 0 {
		handlers[0].Handler(&types.BotContext{})
	}

	handlers = (*dispatchMap)[""][EventMessage]
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
	root := NewRouter("", nil).(*Router)

	route1 := root.OnMessage("test", func(ctx *types.BotContext) {}).(*Route)
	route1.UseState("auth")

	root.OnMessage("help", func(ctx *types.BotContext) {})

	dispatchMap := NewDispatchMap(root)

	if _, exists := (*dispatchMap)["auth"]; !exists {
		t.Error("State 'auth' not found in dispatch map")
	}

	if _, exists := (*dispatchMap)[""]; !exists {
		t.Error("Default state '' not found in dispatch map")
	}

	t.Logf("Dispatch map states: %v", getMapKeys(*dispatchMap))
}

func getMapKeys(m DispatchMap) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func BenchmarkDispatchMapBuild(b *testing.B) {
	root := NewRouter("", nil).(*Router)

	for i := 0; i < 10; i++ {
		group := root.Group("group").(*Router)
		group.Use(func(next types.Handler) types.Handler { return next })

		for j := 0; j < 10; j++ {
			route := group.OnMessage("test", func(ctx *types.BotContext) {}).(*Route)
			route.Use(func(next types.Handler) types.Handler { return next })
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewDispatchMap(root)
	}
}

func BenchmarkDispatchMapAccess(b *testing.B) {
	root := NewRouter("", nil).(*Router)
	root.OnMessage("test", func(ctx *types.BotContext) {})

	dispatchMap := NewDispatchMap(root)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = (*dispatchMap)[""][EventMessage]
	}
}
