package engine

import (
	"github.com/AlexMayka/go-max-sdk/internal/types"
	"github.com/AlexMayka/go-max-sdk/types/models"
	"strings"
)

func (e *botEngine) processUpdate(update *models.Update) types.Handler {
	state := ""
	userID := getUserID(update)
	if userID != 0 {
		if value, ok := e.fsm.GetStateUser(userID); ok {
			state = value
		}
	}

	event := e.getEventType(update)
	if event == types.EventNone {
		return nil
	}

	handler := e.getHandler(state, event, update)

	if handler == nil {
		handler = e.getHandler(state, types.EventAny, update)
	}

	return handler
}

func (e *botEngine) getEventType(update *models.Update) types.EventType {
	eventType, exists := types.UpdateTypeToEvent[update.UpdateType]
	if !exists {
		return types.EventNone
	}

	if eventType == types.EventMessage && update.Message != nil && update.Message.Body != nil {
		text := update.Message.Body.Text
		if len(text) > 0 && text[0] == '/' {
			return types.EventCommand
		}
	}

	return eventType
}

func (e *botEngine) getHandler(state string, event types.EventType, update *models.Update) types.Handler {
	handlers, ok := e.registry.GetHandlers(state, event)
	if !ok {
		return nil
	}

	for _, handler := range handlers {
		msg := getMessage(update)

		switch handler.Match {
		case types.MatchCommand:
			if checkCommand(msg, handler) {
				return handler.Handler
			}
		case types.MatchExact:
			if checkExact(msg, handler) {
				return handler.Handler
			}
		case types.MatchPrefix:
			if checkPrefix(msg, handler) {
				return handler.Handler
			}
		case types.MatchSuffix:
			if checkSuffix(msg, handler) {
				return handler.Handler
			}
		case types.MatchContains:
			if checkContains(msg, handler) {
				return handler.Handler
			}
		case types.MatchRegex:
			if checkRegex(msg, handler) {
				return handler.Handler
			}
		case types.MatchCallback:
			if checkCallback(update.Payload, handler) {
				return handler.Handler
			}
		case types.MatchBotStarted:
			return handler.Handler

		case types.MatchAny:
			return handler.Handler
		}
	}

	return nil
}

func checkCommand(message *string, handler types.RouteHandler) bool {
	if message == nil || handler.Route.GetPattern() == nil {
		return false
	}
	pattern := handler.Route.GetPattern()
	return *pattern == *message
}

func checkCallback(callbackID *string, handler types.RouteHandler) bool {
	if callbackID == nil || handler.Route.GetPattern() == nil {
		return false
	}
	pattern := handler.Route.GetPattern()
	return *pattern == *callbackID
}

func checkExact(message *string, handler types.RouteHandler) bool {
	if message == nil || handler.Route.GetPattern() == nil {
		return false
	}
	pattern := handler.Route.GetPattern()
	return *pattern == *message
}

func checkPrefix(message *string, handler types.RouteHandler) bool {
	if message == nil || handler.Route.GetPattern() == nil {
		return false
	}
	pattern := handler.Route.GetPattern()
	return strings.HasPrefix(*message, *pattern)
}

func checkSuffix(message *string, handler types.RouteHandler) bool {
	pattern := handler.Route.GetPattern()
	if strings.HasSuffix(*message, *pattern) {
		return true
	}

	return false
}

func checkContains(message *string, handler types.RouteHandler) bool {
	if message == nil || handler.Route.GetPattern() == nil {
		return false
	}
	pattern := handler.Route.GetPattern()
	return strings.Contains(*message, *pattern)
}

func checkRegex(message *string, handler types.RouteHandler) bool {
	if message == nil {
		return false
	}

	r := handler.Route.GetCompiledRegex()
	if r == nil {
		return false
	}
	
	return r.MatchString(*message)
}

func getMessage(update *models.Update) *string {
	if update.Message != nil && update.Message.Body != nil {
		return &update.Message.Body.Text
	}
	return nil
}
