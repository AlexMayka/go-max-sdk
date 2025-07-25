package engine

import (
	"context"
	"log"
	"time"

	"github.com/AlexMayka/go-max-sdk/internal/types"
	"github.com/AlexMayka/go-max-sdk/types/models"
)

func (e *botEngine) executeHandlerSimple(handler types.Handler, update *models.Update) {
	userID := getUserID(update)
	chatID := getChatID(update)

	text := ""
	if update.Message != nil && update.Message.Body != nil {
		text = update.Message.Body.Text
	}

	for attempt := 0; attempt <= e.cnf.RetryAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(e.cnf.RetryDelay)
		}

		ctx, cancel := context.WithTimeout(e.ctx, e.cnf.RequestTimeout)

		success := func() bool {
			defer cancel()
			defer func() {
				if r := recover(); r != nil && e.cnf.LogErrors {
					log.Printf("Worker: handler panic (attempt %d): %v", attempt+1, r)
				}
			}()

			botCtx := &types.BotContext{
				Ctx:     ctx,
				Cancel:  cancel,
				UserID:  userID,
				ChatID:  chatID,
				Update:  update,
				Client:  e.client,
				FSM:     e.fsm,
				Text:    text,
				Payload: update.Payload,
			}

			handler(botCtx)
			return true
		}()

		if success {
			return
		}
	}

	if e.cnf.LogErrors {
		log.Printf("Worker: handler failed after %d attempts", e.cnf.RetryAttempts+1)
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
