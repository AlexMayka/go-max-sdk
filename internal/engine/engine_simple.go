package engine

import (
	"context"
	"fmt"
	"github.com/AlexMayka/go-max-sdk/internal/core"
	"github.com/AlexMayka/go-max-sdk/types/models"
)

// executeHandlerSimple executes a handler with retry logic, timeout protection, and panic recovery.
// This is the "simple" synchronous execution model - each handler runs to completion within the worker.
// The handler execution includes:
//   - Retry logic with configurable attempts and delays
//   - Context timeout protection against hanging handlers
//   - Panic recovery with logging
//   - Complete BotContext creation with all update data
func (e *botEngine) executeHandlerSimple(handler core.Handler, update *models.Update) {
	userID := getUserID(update)
	chatID := getChatID(update)
	firstName := getFirstName(update)
	lastName := getLastName(update)
	username := getUsername(update)

	text := ""
	if update.Message != nil && update.Message.Body != nil {
		text = update.Message.Body.Text
	}

	ctx, cancel := context.WithTimeout(e.ctx, e.requestTimeout)

	success := func() bool {
		defer cancel()
		defer func() {
			if r := recover(); r != nil && e.logger != nil {
				e.logger.Error("worker", "handler_panic", fmt.Sprintf("userID=%d, panic=%v", userID, r))
			}
		}()

		botCtx := &core.BotContext{
			Ctx:       ctx,
			Cancel:    cancel,
			UserID:    userID,
			ChatID:    chatID,
			FirstName: firstName,
			LastName:  lastName,
			Username:  username,
			Update:    update,
			Client:    e.client,
			FSM:       e.fsm,
			Text:      text,
		}

		handler(botCtx)
		return true
	}()

	if success {
		return
	}

}

func getChatID(update *models.Update) int64 {
	if update.Message != nil && update.Message.Recipient != nil {
		return update.Message.Recipient.ChatID
	}
	return 0
}

func getUserID(update *models.Update) int64 {
	if update.Message != nil && update.Message.Sender != nil {
		return update.Message.Sender.UserID
	}
	return 0
}

func getFirstName(update *models.Update) string {
	if update.Message != nil && update.Message.Sender != nil {
		return update.Message.Sender.FirstName
	}
	return ""
}

func getLastName(update *models.Update) string {
	if update.Message != nil && update.Message.Sender != nil {
		return update.Message.Sender.LastName
	}
	return ""
}

func getUsername(update *models.Update) string {
	if update.Message != nil && update.Message.Sender != nil {
		return update.Message.Sender.Username
	}
	return ""
}
