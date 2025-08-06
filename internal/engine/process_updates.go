package engine

import (
	"github.com/AlexMayka/go-max-sdk/internal/core"
	"github.com/AlexMayka/go-max-sdk/types/models"
	"strings"
)

// processUpdate processes an incoming update and finds the appropriate handler.
// It determines the user's FSM state, identifies the event type, and matches against registered routes.
// Returns the matched handler or nil if no handler is found.
func (e *botEngine) processUpdate(update *models.Update) core.Handler {
	state := ""
	userID := getUserID(update)
	if userID != 0 {
		if value, ok := e.fsm.GetStateUser(userID); ok {
			state = value
		}
	}

	event := e.getEventType(update)
	if event == core.EventNone {
		return nil
	}

	handler := e.getHandler(state, event, update)

	if handler == nil {
		handler = e.getHandler(state, core.EventAny, update)
	}

	return handler
}

// getEventType determines the event type from an update.
// It converts update types to event types and handles special cases like commands.
func (e *botEngine) getEventType(update *models.Update) core.EventType {
	eventType, exists := core.UpdateTypeToEvent[update.UpdateType]
	if !exists {
		return core.EventNone
	}

	if eventType == core.EventMessage && update.Message != nil && update.Message.Body != nil {
		text := update.Message.Body.Text
		if len(text) > 0 && text[0] == '/' {
			return core.EventCommand
		}
	}

	return eventType
}

// getHandler finds a matching handler for the given state, event type, and update.
// It iterates through registered handlers and applies pattern matching based on the handler's match type.
func (e *botEngine) getHandler(state string, event core.EventType, update *models.Update) core.Handler {
	handlers, ok := e.registry.GetHandlers(state, event)
	if !ok {
		return nil
	}

	for _, handler := range handlers {
		msg := getMessage(update)

		switch handler.Match {
		case core.MatchCommand:
			if checkCommand(msg, handler) {
				return handler.Handler
			}
		case core.MatchExact:
			if checkExact(msg, handler) {
				return handler.Handler
			}
		case core.MatchPrefix:
			if checkPrefix(msg, handler) {
				return handler.Handler
			}
		case core.MatchSuffix:
			if checkSuffix(msg, handler) {
				return handler.Handler
			}
		case core.MatchContains:
			if checkContains(msg, handler) {
				return handler.Handler
			}
		case core.MatchRegex:
			if checkRegex(msg, handler) {
				return handler.Handler
			}
		case core.MatchCallback:
			if checkCallback(&update.Message.Body.Text, handler) {
				return handler.Handler
			}
		case core.MatchBotStarted:
			return handler.Handler

		case core.MatchAny:
			return handler.Handler
		}
	}

	return nil
}

// checkCommand verifies if a message matches a command pattern.
// Commands are exact matches with the registered command pattern.
func checkCommand(message *string, handler core.RouteHandler) bool {
	if message == nil || handler.Route.GetPattern() == nil {
		return false
	}
	pattern := handler.Route.GetPattern()
	return *pattern == *message
}

// checkCallback verifies if a callback ID matches the registered callback pattern.
func checkCallback(callbackID *string, handler core.RouteHandler) bool {
	if callbackID == nil || handler.Route.GetPattern() == nil {
		return false
	}
	pattern := handler.Route.GetPattern()
	return *pattern == *callbackID
}

// checkExact verifies if a message exactly matches the registered pattern.
func checkExact(message *string, handler core.RouteHandler) bool {
	if message == nil || handler.Route.GetPattern() == nil {
		return false
	}
	pattern := handler.Route.GetPattern()
	return *pattern == *message
}

// checkPrefix verifies if a message starts with the registered prefix pattern.
func checkPrefix(message *string, handler core.RouteHandler) bool {
	if message == nil || handler.Route.GetPattern() == nil {
		return false
	}
	pattern := handler.Route.GetPattern()
	return strings.HasPrefix(*message, *pattern)
}

// checkSuffix verifies if a message ends with the registered suffix pattern.
func checkSuffix(message *string, handler core.RouteHandler) bool {
	pattern := handler.Route.GetPattern()
	if strings.HasSuffix(*message, *pattern) {
		return true
	}

	return false
}

// checkContains verifies if a message contains the registered substring pattern.
func checkContains(message *string, handler core.RouteHandler) bool {
	if message == nil || handler.Route.GetPattern() == nil {
		return false
	}
	pattern := handler.Route.GetPattern()
	return strings.Contains(*message, *pattern)
}

// checkRegex verifies if a message matches the registered regular expression pattern.
func checkRegex(message *string, handler core.RouteHandler) bool {
	if message == nil {
		return false
	}

	r := handler.Route.GetCompiledRegex()
	if r == nil {
		return false
	}

	return r.MatchString(*message)
}

// getMessage extracts the text message from an update.
// Returns nil if the update doesn't contain a text message.
func getMessage(update *models.Update) *string {
	if update.Message != nil && update.Message.Body != nil {
		return &update.Message.Body.Text
	}
	return nil
}
