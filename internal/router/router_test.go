package router

import (
	"github.com/AlexMayka/go-max-sdk/internal/core"
	"regexp"
	"testing"
)

func TestRouterBasicOperations(t *testing.T) {
	t.Run("CreateRouter", func(t *testing.T) {
		router := NewRouter("test", nil)
		if router == nil {
			t.Error("Expected router to be created")
		}

		r := router.(*Router)
		if r.Prefix != "test" {
			t.Errorf("Expected prefix 'test', got '%s'", r.Prefix)
		}
		if r.Parent != nil {
			t.Error("Expected parent to be nil")
		}
	})

	t.Run("DefaultRouter", func(t *testing.T) {
		router := DefaultRouter()
		if router == nil {
			t.Error("Expected default router to be created")
		}

		r := router.(*Router)
		if r.Prefix != "" {
			t.Errorf("Expected empty prefix, got '%s'", r.Prefix)
		}
	})

	t.Run("CreateGroup", func(t *testing.T) {
		root := NewRouter("", nil).(*Router)
		group := root.Group("api")

		if group == nil {
			t.Error("Expected group to be created")
		}

		g := group.(*Router)
		if g.Prefix != "api" {
			t.Errorf("Expected group prefix 'api', got '%s'", g.Prefix)
		}
		if g.Parent != root {
			t.Error("Expected group parent to be root router")
		}

		if len(root.Children) != 1 {
			t.Errorf("Expected root to have 1 child, got %d", len(root.Children))
		}
		if root.Children[0] != g {
			t.Error("Expected group to be in parent's children")
		}
	})
}

func TestRouterMiddleware(t *testing.T) {
	middleware1 := func(next core.Handler) core.Handler {
		return func(ctx *core.BotContext) {
			next(ctx)
		}
	}
	middleware2 := func(next core.Handler) core.Handler {
		return func(ctx *core.BotContext) {
			next(ctx)
		}
	}

	t.Run("UseMiddleware", func(t *testing.T) {
		router := NewRouter("test", nil).(*Router)

		result := router.Use(middleware1, middleware2)

		if result != router {
			t.Error("Expected Use to return the same router for chaining")
		}

		if len(router.Middlewares) != 2 {
			t.Errorf("Expected 2 middlewares, got %d", len(router.Middlewares))
		}
	})

	t.Run("UseState", func(t *testing.T) {
		router := NewRouter("test", nil).(*Router)

		result := router.UseState("waiting_input")

		if result != router {
			t.Error("Expected UseState to return the same router for chaining")
		}

		if router.State != "waiting_input" {
			t.Errorf("Expected state 'waiting_input', got '%s'", router.State)
		}
	})
}

func TestRouterRoutes(t *testing.T) {
	handler := func(ctx *core.BotContext) {}

	t.Run("OnStarted", func(t *testing.T) {
		router := NewRouter("", nil).(*Router)
		route := router.OnStarted(handler)

		if route == nil {
			t.Error("Expected route to be created")
		}

		if len(router.Routes) != 1 {
			t.Errorf("Expected 1 route, got %d", len(router.Routes))
		}

		r := router.Routes[0]
		if r.GetType() != core.EventBotStarted {
			t.Errorf("Expected route type RouteBotStarted, got %v", r.GetType())
		}
		if r.GetPattern() != nil {
			t.Error("Expected pattern to be nil for OnStarted")
		}
	})

	t.Run("OnMessage", func(t *testing.T) {
		router := NewRouter("", nil).(*Router)
		route := router.OnMessage("hello", handler)

		if route == nil {
			t.Error("Expected route to be created")
		}

		r := router.Routes[0]
		if r.GetType() != core.EventMessage {
			t.Errorf("Expected route type RouteMsg, got %v", r.GetType())
		}
		if r.GetPattern() == nil || *r.GetPattern() != "hello" {
			t.Errorf("Expected pattern 'hello', got %v", r.GetPattern())
		}
		if r.GetPrefix() != ":hello" {
			t.Errorf("Expected prefix ':hello', got '%s'", r.GetPrefix())
		}
	})

	t.Run("OnCommand", func(t *testing.T) {
		router := NewRouter("", nil).(*Router)
		route := router.OnCommand("start", handler)

		if route == nil {
			t.Error("Expected route to be created")
		}

		r := router.Routes[0]
		if r.GetType() != core.EventCommand {
			t.Errorf("Expected route type RouteCommand, got %v", r.GetType())
		}
		if r.GetPattern() == nil || *r.GetPattern() != "start" {
			t.Errorf("Expected pattern 'start', got %v", r.GetPattern())
		}
	})

	t.Run("OnCallback", func(t *testing.T) {
		router := NewRouter("", nil).(*Router)
		route := router.OnCallback("button_click", handler)

		if route == nil {
			t.Error("Expected route to be created")
		}

		r := router.Routes[0]
		if r.GetType() != core.EventCallback {
			t.Errorf("Expected route type RouteCallback, got %v", r.GetType())
		}
		if r.GetPattern() == nil || *r.GetPattern() != "button_click" {
			t.Errorf("Expected pattern 'button_click', got %v", r.GetPattern())
		}
	})

	t.Run("Any", func(t *testing.T) {
		router := NewRouter("", nil).(*Router)
		route := router.Any(handler)

		if route == nil {
			t.Error("Expected route to be created")
		}

		r := router.GetAnyMsg()
		if r.GetType() != core.EventAny {
			t.Errorf("Expected route type RouteAny, got %v", r.GetType())
		}
		if r.GetPattern() != nil {
			t.Error("Expected pattern to be nil for Any")
		}
	})
}

func TestRouterRegexValidation(t *testing.T) {
	handler := func(ctx *core.BotContext) {}

	t.Run("ValidRegex", func(t *testing.T) {
		router := NewRouter("", nil).(*Router)

		route := router.OnRegex(`^\d+$`, handler)

		if route == nil {
			t.Error("Expected route to be created")
		}

		if len(router.Routes) == 0 {
			t.Fatal("Expected at least one route to be created")
		}

		r := router.Routes[0]
		if r.GetType() != core.EventMessage {
			t.Errorf("Expected route type EventMessage, got %v", r.GetType())
		}
		if r.GetPattern() == nil || *r.GetPattern() != `^\d+$` {
			t.Errorf("Expected pattern '^\\d+$', got %v", r.GetPattern())
		}
	})

	t.Run("InvalidRegexPanics", func(t *testing.T) {
		router := NewRouter("", nil).(*Router)

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for invalid regex")
			} else {
				panicMsg := r.(string)
				if len(panicMsg) == 0 {
					t.Error("Expected non-empty panic message")
				}
			}
		}()

		router.OnRegex(`[unclosed`, handler)
	})

	t.Run("ComplexValidRegex", func(t *testing.T) {
		router := NewRouter("", nil).(*Router)

		validPatterns := []string{
			`^/start`,
			`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`,
			`[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z]{2,}`,
			`^(https?|ftp)://[^\s/$.?#].[^\s]*$`,
		}

		for _, pattern := range validPatterns {
			route := router.OnRegex(pattern, handler)
			if route == nil {
				t.Errorf("Expected route to be created for pattern: %s", pattern)
			}
		}

		if len(router.Routes) != len(validPatterns) {
			t.Errorf("Expected %d routes, got %d", len(validPatterns), len(router.Routes))
		}
	})
}

func TestRouterPrefixBuilding(t *testing.T) {
	handler := func(ctx *core.BotContext) {}

	t.Run("PrefixWithPattern", func(t *testing.T) {
		router := NewRouter("api", nil).(*Router)
		_ = router.OnMessage("users", handler)

		r := router.Routes[0]
		expectedPrefix := "api:users"
		if r.GetPrefix() != expectedPrefix {
			t.Errorf("Expected prefix '%s', got '%s'", expectedPrefix, r.GetPrefix())
		}
	})

	t.Run("PrefixWithoutPattern", func(t *testing.T) {
		router := NewRouter("api", nil).(*Router)
		_ = router.OnStarted(handler)

		r := router.Routes[0]
		expectedPrefix := "api"
		if r.GetPrefix() != expectedPrefix {
			t.Errorf("Expected prefix '%s', got '%s'", expectedPrefix, r.GetPrefix())
		}
	})

	t.Run("EmptyPrefixWithPattern", func(t *testing.T) {
		router := NewRouter("", nil).(*Router)
		_ = router.OnCommand("help", handler)

		r := router.Routes[0]
		expectedPrefix := ":help"
		if r.GetPrefix() != expectedPrefix {
			t.Errorf("Expected prefix '%s', got '%s'", expectedPrefix, r.GetPrefix())
		}
	})

	t.Run("EmptyPrefixWithoutPattern", func(t *testing.T) {
		router := NewRouter("", nil).(*Router)
		_ = router.Any(handler)

		r := router.GetAnyMsg()
		expectedPrefix := ""
		if r.GetPrefix() != expectedPrefix {
			t.Errorf("Expected prefix '%s', got '%s'", expectedPrefix, r.GetPrefix())
		}
	})
}

func TestRouterChaining(t *testing.T) {
	handler := func(ctx *core.BotContext) {}
	middleware := func(next core.Handler) core.Handler {
		return func(ctx *core.BotContext) { next(ctx) }
	}

	t.Run("FluentAPI", func(t *testing.T) {
		router := NewRouter("", nil)

		result := router.Use(middleware).UseState("active").Group("api").Use(middleware)

		if result == nil {
			t.Error("Expected chaining to work")
		}

		group := result.(*Router)
		if group.Prefix != "api" {
			t.Error("Expected group to be created through chaining")
		}
		if len(group.Middlewares) != 1 {
			t.Error("Expected middleware to be added through chaining")
		}
	})

	t.Run("RouteChaining", func(t *testing.T) {
		router := NewRouter("", nil)

		route := router.OnMessage("test", handler).Use(middleware).UseState("test_state")

		if route == nil {
			t.Error("Expected route chaining to work")
		}

		r := router.(*Router).Routes[0]
		if len(r.GetMiddleware()) != 1 {
			t.Error("Expected middleware to be added to route")
		}
		if r.GetState() != "test_state" {
			t.Error("Expected state to be set on route")
		}
	})
}

func TestRouterEdgeCases(t *testing.T) {
	handler := func(ctx *core.BotContext) {}

	t.Run("EmptyPatterns", func(t *testing.T) {
		router := NewRouter("", nil).(*Router)

		router.OnMessage("", handler)
		router.OnCommand("", handler)
		router.OnCallback("", handler)

		if len(router.Routes) != 3 {
			t.Errorf("Expected 3 routes with empty patterns, got %d", len(router.Routes))
		}

		for _, route := range router.Routes {
			if route.GetPattern() == nil || *route.GetPattern() != "" {
				t.Error("Expected empty string pattern to be preserved")
			}
		}
	})

	t.Run("SpecialCharacterPatterns", func(t *testing.T) {
		router := NewRouter("", nil).(*Router)

		specialPatterns := []string{
			"hello world",
			"test@example.com",
			"🎉 emoji test",
			"test\nwith\nnewlines",
			"test\twith\ttabs",
		}

		for _, pattern := range specialPatterns {
			router.OnMessage(pattern, handler)
		}

		if len(router.Routes) != len(specialPatterns) {
			t.Errorf("Expected %d routes, got %d", len(specialPatterns), len(router.Routes))
		}

		for i, route := range router.Routes {
			if route.GetPattern() == nil || *route.GetPattern() != specialPatterns[i] {
				t.Errorf("Expected pattern '%s', got %v", specialPatterns[i], route.GetPattern())
			}
		}
	})

	t.Run("MultipleRoutesOfSameType", func(t *testing.T) {
		router := NewRouter("", nil).(*Router)

		router.OnMessage("hello", handler)
		router.OnMessage("hi", handler)
		router.OnMessage("hey", handler)

		if len(router.Routes) != 3 {
			t.Errorf("Expected 3 routes, got %d", len(router.Routes))
		}

		for _, route := range router.Routes {
			if route.GetType() != core.EventMessage {
				t.Error("Expected all routes to be EventMessage type")
			}
		}
	})

	t.Run("NilHandler", func(t *testing.T) {
		router := NewRouter("", nil).(*Router)

		route := router.OnMessage("test", nil)

		if route == nil {
			t.Error("Expected route to be created even with nil handler")
		}

		r := router.Routes[0]
		if r.GetHandler() != nil {
			t.Error("Expected handler to be nil")
		}
	})
}

func TestRegexCompilation(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		valid   bool
	}{
		{"Simple pattern", "hello", true},
		{"Digit pattern", `\d+`, true},
		{"Email pattern", `[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z]{2,}`, true},
		{"Start anchor", "^start", true},
		{"End anchor", "end$", true},
		{"Both anchors", "^exact$", true},
		{"Invalid bracket", "[unclosed", false},
		{"Invalid group", "(unclosed", false},
		{"Invalid escape", `\z`, true},
		{"Invalid quantifier", "*invalid", false},
	}

	handler := func(ctx *core.BotContext) {}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := NewRouter("", nil).(*Router)

			_, err := regexp.Compile(tc.pattern)

			if tc.valid && err != nil {
				t.Errorf("Expected pattern '%s' to be valid, but got error: %v", tc.pattern, err)
			}

			if !tc.valid && err == nil {
				t.Errorf("Expected pattern '%s' to be invalid, but it compiled successfully", tc.pattern)
			}

			if tc.valid {
				route := router.OnRegex(tc.pattern, handler)
				if route == nil {
					t.Error("Expected route to be created for valid pattern")
				}
			} else {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic for invalid regex pattern")
					}
				}()
				router.OnRegex(tc.pattern, handler)
			}
		})
	}
}
