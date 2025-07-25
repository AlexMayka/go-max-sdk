package engine

import (
	"regexp"
	"testing"

	"github.com/AlexMayka/go-max-sdk/internal/types"
	"github.com/AlexMayka/go-max-sdk/types/models"
)

type mockFSM struct {
	states map[int64]string
	values map[int64]map[string]string
}

func (m *mockFSM) CreateUser(id int64, state string) bool {
	if m.states == nil {
		m.states = make(map[int64]string)
	}
	m.states[id] = state
	return true
}

func (m *mockFSM) SetStateUser(userID int64, state string) bool {
	if m.states == nil {
		m.states = make(map[int64]string)
	}
	m.states[userID] = state
	return true
}

func (m *mockFSM) GetStateUser(userID int64) (string, bool) {
	if m.states == nil {
		return "", false
	}
	state, exists := m.states[userID]
	return state, exists
}

func (m *mockFSM) SetValue(id int64, param string, data string) bool {
	if m.values == nil {
		m.values = make(map[int64]map[string]string)
	}
	if m.values[id] == nil {
		m.values[id] = make(map[string]string)
	}
	m.values[id][param] = data
	return true
}

func (m *mockFSM) GetValue(id int64, param string) (string, bool) {
	if m.values == nil || m.values[id] == nil {
		return "", false
	}
	value, exists := m.values[id][param]
	return value, exists
}

type mockRegistry struct {
	handlers map[string]map[types.EventType][]types.RouteHandler
	states   []string
}

func (m *mockRegistry) GetHandlers(state string, event types.EventType) ([]types.RouteHandler, bool) {
	if m.handlers == nil {
		return nil, false
	}
	stateHandlers, exists := m.handlers[state]
	if !exists {
		return nil, false
	}
	handlers, exists := stateHandlers[event]
	return handlers, exists
}

func (m *mockRegistry) HasState(state string) bool {
	for _, s := range m.states {
		if s == state {
			return true
		}
	}
	return false
}

func (m *mockRegistry) GetStates() []string {
	return m.states
}

func (m *mockRegistry) AddHandler(state string, event types.EventType, handler types.RouteHandler) {
	if m.handlers == nil {
		m.handlers = make(map[string]map[types.EventType][]types.RouteHandler)
	}
	if m.handlers[state] == nil {
		m.handlers[state] = make(map[types.EventType][]types.RouteHandler)
	}
	m.handlers[state][event] = append(m.handlers[state][event], handler)

	if !m.HasState(state) {
		m.states = append(m.states, state)
	}
}

type mockRoute struct {
	pattern       *string
	handler       types.Handler
	matchType     types.MatchType
	eventType     types.EventType
	prefix        string
	state         string
	middleware    []types.Middleware
	compiledRegex *regexp.Regexp
}

func (m *mockRoute) GetPattern() *string {
	return m.pattern
}

func (m *mockRoute) GetHandler() types.Handler {
	return m.handler
}

func (m *mockRoute) GetMatchType() types.MatchType {
	return m.matchType
}

func (m *mockRoute) GetType() types.EventType {
	return m.eventType
}

func (m *mockRoute) GetPrefix() string {
	return m.prefix
}

func (m *mockRoute) GetState() string {
	return m.state
}

func (m *mockRoute) GetMiddleware() []types.Middleware {
	return m.middleware
}

func (m *mockRoute) GetCompiledRegex() *regexp.Regexp {
	return m.compiledRegex
}

func (m *mockRoute) UseState(state string) types.Route {
	m.state = state
	return m
}

func (m *mockRoute) Use(middlewares ...types.Middleware) types.Route {
	m.middleware = append(m.middleware, middlewares...)
	return m
}

func newMockRoute(pattern string, handler types.Handler, matchType types.MatchType, eventType types.EventType) *mockRoute {
	var compiledRegex *regexp.Regexp
	if matchType == types.MatchRegex && pattern != "" {
		var err error
		compiledRegex, err = regexp.Compile(pattern)
		if err != nil {
			panic(err)
		}
	}

	return &mockRoute{
		pattern:       &pattern,
		handler:       handler,
		matchType:     matchType,
		eventType:     eventType,
		compiledRegex: compiledRegex,
	}
}

func newMockRouteHandler(match types.MatchType, pattern string, handler types.Handler, eventType types.EventType) types.RouteHandler {
	return types.RouteHandler{
		Match:   match,
		Handler: handler,
		Route:   newMockRoute(pattern, handler, match, eventType),
	}
}

func TestProcessUpdate(t *testing.T) {
	mockHandler := func(ctx *types.BotContext) {}

	t.Run("MessageWithoutState", func(t *testing.T) {
		fsm := &mockFSM{}
		registry := &mockRegistry{}

		// Добавляем handler для пустого state
		registry.AddHandler("", types.EventMessage, newMockRouteHandler(
			types.MatchExact, "hello", mockHandler, types.EventMessage,
		))

		engine := &botEngine{
			registry: registry,
			fsm:      fsm,
		}

		update := &models.Update{
			UpdateType: models.UpdateTypeMessageCreated,
			Message: &models.Message{
				Sender: &models.User{UserID: 123},
				Body:   &models.MessageBody{Text: "hello"},
			},
		}

		handler := engine.processUpdate(update)
		if handler == nil {
			t.Error("Expected handler to be found")
		}
	})

	t.Run("MessageWithState", func(t *testing.T) {
		fsm := &mockFSM{}
		registry := &mockRegistry{}

		// Устанавливаем состояние пользователя
		fsm.SetStateUser(123, "waiting_input")

		// Добавляем handler для конкретного state
		registry.AddHandler("waiting_input", types.EventMessage, newMockRouteHandler(
			types.MatchExact, "response", mockHandler, types.EventMessage,
		))

		engine := &botEngine{
			registry: registry,
			fsm:      fsm,
		}

		update := &models.Update{
			UpdateType: models.UpdateTypeMessageCreated,
			Message: &models.Message{
				Sender: &models.User{UserID: 123},
				Body:   &models.MessageBody{Text: "response"},
			},
		}

		handler := engine.processUpdate(update)
		if handler == nil {
			t.Error("Expected handler to be found for state")
		}
	})

	t.Run("CommandDetection", func(t *testing.T) {
		fsm := &mockFSM{}
		registry := &mockRegistry{}

		registry.AddHandler("", types.EventCommand, newMockRouteHandler(
			types.MatchCommand, "/start", mockHandler, types.EventCommand,
		))

		engine := &botEngine{
			registry: registry,
			fsm:      fsm,
		}

		update := &models.Update{
			UpdateType: models.UpdateTypeMessageCreated,
			Message: &models.Message{
				Sender: &models.User{UserID: 123},
				Body:   &models.MessageBody{Text: "/start"},
			},
		}

		handler := engine.processUpdate(update)
		if handler == nil {
			t.Error("Expected command handler to be found")
		}
	})

	t.Run("FallbackToAny", func(t *testing.T) {
		fsm := &mockFSM{}
		registry := &mockRegistry{}

		// Добавляем только Any handler
		registry.AddHandler("", types.EventAny, newMockRouteHandler(
			types.MatchAny, "", mockHandler, types.EventAny,
		))

		engine := &botEngine{
			registry: registry,
			fsm:      fsm,
		}

		update := &models.Update{
			UpdateType: models.UpdateTypeMessageCreated,
			Message: &models.Message{
				Sender: &models.User{UserID: 123},
				Body:   &models.MessageBody{Text: "anything"},
			},
		}

		handler := engine.processUpdate(update)
		if handler == nil {
			t.Error("Expected Any handler to be found as fallback")
		}
	})

	t.Run("NoHandlerFound", func(t *testing.T) {
		fsm := &mockFSM{}
		registry := &mockRegistry{}

		engine := &botEngine{
			registry: registry,
			fsm:      fsm,
		}

		update := &models.Update{
			UpdateType: models.UpdateTypeMessageCreated,
			Message: &models.Message{
				Sender: &models.User{UserID: 123},
				Body:   &models.MessageBody{Text: "no match"},
			},
		}

		handler := engine.processUpdate(update)
		if handler != nil {
			t.Error("Expected no handler to be found")
		}
	})
}

func TestGetEventType(t *testing.T) {
	engine := &botEngine{}

	t.Run("MessageEvent", func(t *testing.T) {
		update := &models.Update{
			UpdateType: models.UpdateTypeMessageCreated,
			Message: &models.Message{
				Body: &models.MessageBody{Text: "hello"},
			},
		}

		eventType := engine.getEventType(update)
		if eventType != types.EventMessage {
			t.Errorf("Expected EventMessage, got %v", eventType)
		}
	})

	t.Run("CommandEvent", func(t *testing.T) {
		update := &models.Update{
			UpdateType: models.UpdateTypeMessageCreated,
			Message: &models.Message{
				Body: &models.MessageBody{Text: "/start"},
			},
		}

		eventType := engine.getEventType(update)
		if eventType != types.EventCommand {
			t.Errorf("Expected EventCommand, got %v", eventType)
		}
	})

	t.Run("CallbackEvent", func(t *testing.T) {
		update := &models.Update{
			UpdateType: models.UpdateTypeMessageCallback,
		}

		eventType := engine.getEventType(update)
		if eventType != types.EventCallback {
			t.Errorf("Expected EventCallback, got %v", eventType)
		}
	})

	t.Run("UnknownEvent", func(t *testing.T) {
		update := &models.Update{
			UpdateType: "unknown_type",
		}

		eventType := engine.getEventType(update)
		if eventType != types.EventNone {
			t.Errorf("Expected EventNone, got %v", eventType)
		}
	})
}

func TestMatching(t *testing.T) {
	t.Run("CheckCommand", func(t *testing.T) {
		pattern := "/start"
		message := "/start"
		handler := newMockRouteHandler(types.MatchCommand, pattern, nil, types.EventCommand)

		result := checkCommand(&message, handler)
		if !result {
			t.Error("Expected command to match")
		}

		wrongMessage := "/help"
		result = checkCommand(&wrongMessage, handler)
		if result {
			t.Error("Expected command not to match")
		}
	})

	t.Run("CheckExact", func(t *testing.T) {
		pattern := "hello"
		message := "hello"
		handler := newMockRouteHandler(types.MatchExact, pattern, nil, types.EventMessage)

		result := checkExact(&message, handler)
		if !result {
			t.Error("Expected exact match")
		}

		wrongMessage := "Hello"
		result = checkExact(&wrongMessage, handler)
		if result {
			t.Error("Expected exact match to be case sensitive")
		}
	})

	t.Run("CheckPrefix", func(t *testing.T) {
		pattern := "hello"
		message := "hello world"
		handler := newMockRouteHandler(types.MatchPrefix, pattern, nil, types.EventMessage)

		result := checkPrefix(&message, handler)
		if !result {
			t.Error("Expected prefix match")
		}

		wrongMessage := "hi world"
		result = checkPrefix(&wrongMessage, handler)
		if result {
			t.Error("Expected prefix not to match")
		}
	})

	t.Run("CheckSuffix", func(t *testing.T) {
		pattern := "world"
		message := "hello world"
		handler := newMockRouteHandler(types.MatchSuffix, pattern, nil, types.EventMessage)

		result := checkSuffix(&message, handler)
		if !result {
			t.Error("Expected suffix match")
		}

		wrongMessage := "hello there"
		result = checkSuffix(&wrongMessage, handler)
		if result {
			t.Error("Expected suffix not to match")
		}
	})

	t.Run("CheckContains", func(t *testing.T) {
		pattern := "test"
		message := "this is a test message"
		handler := newMockRouteHandler(types.MatchContains, pattern, nil, types.EventMessage)

		result := checkContains(&message, handler)
		if !result {
			t.Error("Expected contains match")
		}

		wrongMessage := "this is a message"
		result = checkContains(&wrongMessage, handler)
		if result {
			t.Error("Expected contains not to match")
		}
	})

	t.Run("CheckRegex", func(t *testing.T) {
		pattern := `^\d+$`
		message := "12345"
		handler := newMockRouteHandler(types.MatchRegex, pattern, nil, types.EventMessage)

		result := checkRegex(&message, handler)
		if !result {
			t.Error("Expected regex match")
		}

		wrongMessage := "abc123"
		result = checkRegex(&wrongMessage, handler)
		if result {
			t.Error("Expected regex not to match")
		}

		// Test nil regex (when compilation failed)
		nilRegexHandler := newMockRouteHandler(types.MatchExact, "valid", nil, types.EventMessage) // не regex
		nilRegexHandler.Route.(*mockRoute).compiledRegex = nil
		result = checkRegex(&message, nilRegexHandler)
		if result {
			t.Error("Expected nil regex not to match")
		}
	})

	t.Run("CheckCallback", func(t *testing.T) {
		pattern := "button_1"
		callbackID := "button_1"
		handler := newMockRouteHandler(types.MatchCallback, pattern, nil, types.EventCallback)

		result := checkCallback(&callbackID, handler)
		if !result {
			t.Error("Expected callback match")
		}

		wrongCallback := "button_2"
		result = checkCallback(&wrongCallback, handler)
		if result {
			t.Error("Expected callback not to match")
		}
	})

	t.Run("NilInputs", func(t *testing.T) {
		pattern := "test"
		handler := newMockRouteHandler(types.MatchExact, pattern, nil, types.EventMessage)

		if checkExact(nil, handler) {
			t.Error("Expected nil message not to match")
		}

		if checkPrefix(nil, handler) {
			t.Error("Expected nil message not to match")
		}

		if checkContains(nil, handler) {
			t.Error("Expected nil message not to match")
		}

		if checkRegex(nil, handler) {
			t.Error("Expected nil message not to match")
		}

		if checkCallback(nil, handler) {
			t.Error("Expected nil callback not to match")
		}
	})
}

func TestGetMessage(t *testing.T) {
	t.Run("ValidMessage", func(t *testing.T) {
		update := &models.Update{
			Message: &models.Message{
				Body: &models.MessageBody{Text: "hello"},
			},
		}

		message := getMessage(update)
		if message == nil || *message != "hello" {
			t.Errorf("Expected 'hello', got %v", message)
		}
	})

	t.Run("NoMessage", func(t *testing.T) {
		update := &models.Update{}

		message := getMessage(update)
		if message != nil {
			t.Error("Expected nil message")
		}
	})

	t.Run("NoBody", func(t *testing.T) {
		update := &models.Update{
			Message: &models.Message{},
		}

		message := getMessage(update)
		if message != nil {
			t.Error("Expected nil message when body is nil")
		}
	})
}
