package maxsdk

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexMayka/go-max-sdk/internal"
	"github.com/AlexMayka/go-max-sdk/internal/api"
	"github.com/AlexMayka/go-max-sdk/internal/core"
	"github.com/AlexMayka/go-max-sdk/internal/engine"
	"github.com/AlexMayka/go-max-sdk/internal/fsm/local"
	"github.com/AlexMayka/go-max-sdk/internal/registry"
	"github.com/AlexMayka/go-max-sdk/internal/router"
	"github.com/AlexMayka/go-max-sdk/internal/transport"
	"github.com/AlexMayka/go-max-sdk/middleware"
	"github.com/AlexMayka/go-max-sdk/types/models"
	reqMsg "github.com/AlexMayka/go-max-sdk/types/requests/messages"
)

// Type aliases for public API

// BotContext provides access to the current update and bot functionality within handlers.
type BotContext = core.BotContext

// Handler is a function that processes incoming updates.
type Handler = core.Handler

// Middleware is a function that wraps handlers to provide cross-cutting functionality.
type Middleware = core.Middleware

// Router manages message routing and middleware application.
type Router = core.Router

// Route represents a registered message handler with its matching criteria.
type Route = core.Route

// Config contains bot configuration settings like timeouts and limits.
type Config = internal.Config

// FSM manages finite state machine for conversation flows.
type FSM = core.FSM

// Logger provides structured logging interface for the bot.
type Logger = core.Logger

// Transport handles communication with the MAX API (polling, webhooks, etc.).
type Transport = core.Transport

// DefaultRouter creates a new default router instance.
func DefaultRouter() Router {
	return router.DefaultRouter()
}

// Bot represents a MAX bot instance with routing, middleware, and messaging capabilities.
type Bot struct {
	token string

	// Public configuration fields
	Config    *Config   // Bot configuration settings
	Fsm       FSM       // Finite state machine for conversation flows
	Logger    Logger    // Logger interface for structured logging
	Transport Transport // Transport layer (polling, webhook, etc.)

	Logging bool // Enable/disable built-in logging

	// Internal components
	client   core.APIClient
	registry core.RouteRegistry
	engine   core.BotEngine

	router Router
}

// NewBot creates a new Bot instance with the provided token.
// Initializes default configuration, FSM, logger, and router.
func NewBot(token string) *Bot {
	return &Bot{
		token: token,

		Config:  internal.DefaultConfig(),
		Fsm:     local.NewFSM(),
		Logger:  internal.NewConsoleLogger(),
		Logging: true,

		router: router.DefaultRouter(),
	}
}

// OnStarted registers a handler for the bot start event.
func (b *Bot) OnStarted(handler internal.Handler) Route {
	return b.router.OnStarted(handler)
}

// OnMessage registers a handler for exact message text matches.
func (b *Bot) OnMessage(msg string, handler internal.Handler) Route {
	return b.router.OnMessage(msg, handler)
}

// OnRegex registers a handler for messages matching the given regular expression.
func (b *Bot) OnRegex(regex string, handler internal.Handler) Route {
	return b.router.OnRegex(regex, handler)
}

// OnPrefix registers a handler for messages starting with the specified prefix.
func (b *Bot) OnPrefix(prefix string, handler internal.Handler) Route {
	return b.router.OnPrefix(prefix, handler)
}

// OnSuffix registers a handler for messages ending with the specified suffix.
func (b *Bot) OnSuffix(suffix string, handler internal.Handler) Route {
	return b.router.OnSuffix(suffix, handler)
}

// OnContains registers a handler for messages containing the specified substring.
func (b *Bot) OnContains(sub string, handler internal.Handler) Route {
	return b.router.OnContains(sub, handler)
}

// OnCommand registers a handler for bot commands (e.g., "/start", "/help").
func (b *Bot) OnCommand(cmd string, handler internal.Handler) Route {
	return b.router.OnCommand(cmd, handler)
}

// OnCallback registers a handler for inline keyboard callback data.
func (b *Bot) OnCallback(call string, handler internal.Handler) Route {
	return b.router.OnCallback(call, handler)
}

// Any registers a handler that matches any message (catch-all).
func (b *Bot) Any(handler internal.Handler) Route {
	return b.router.Any(handler)
}

// Group creates a new router group with the specified prefix.
func (b *Bot) Group(prefix string) Router {
	return b.router.Group(prefix)
}

// Use applies middleware to all routes registered after this call.
func (b *Bot) Use(middlewares ...internal.Middleware) Router {
	return b.router.Use(middlewares...)
}

// UseState creates a router that only matches messages from users in the specified FSM state.
func (b *Bot) UseState(state string) Router {
	return b.router.UseState(state)
}

// SetRouter replaces the default router with a custom implementation.
func (b *Bot) SetRouter(router Router) {
	b.router = router
}

// Start initializes and starts the bot. This method blocks until the bot is stopped.
func (b *Bot) Start() error {
	b.client = api.NewClient(
		b.token,
		b.Logger,
		b.Config.APITimeout,
		b.Config.APIBurstLimit,
		b.Config.APIRateLimit,
		b.Config.APIMaxRetries,
		b.Config.APIRetryDelay,
	)

	if b.registry == nil {
		b.registry = registry.BuildFrom(b.router)
	}

	b.router.Use(middleware.Logging(b.Logger))

	if b.Transport == nil {
		b.Transport = transport.NewLongPolling(
			b.Config.PollTimeout,
			b.Config.PollLimit,
			b.Config.UpdatesBuffer,
			b.Config.PollRetryDelay,
			b.Config.PollMaxRetries,
			b.client,
			b.Logger,
		)
	}

	if b.engine == nil {
		b.engine = engine.NewBotEngine(
			b.registry,
			b.Fsm,
			b.client,
			b.Transport,
			context.Background(),
			b.Config.MaxWorkers,
			b.Config.WorkerQueueSize,
			b.Config.RequestTimeout,
			b.Config.ShutdownTimeout,
			b.Config.RetryAttempts,
			b.Config.RetryDelay,
			b.Logger,
		)
	}

	return b.engine.Start()
}

// Stop gracefully shuts down the bot.
func (b *Bot) Stop() error {
	if b.engine == nil {
		return errors.New("bot not started")
	}
	return b.engine.Stop()
}

// SendText sends a text message to a user or chat.
// Provide either userID (for private message) or chatID (for group message), not both.
func (b *Bot) SendText(userID, chatID int64, text string) error {
	body := &models.NewMessageBody{Text: &text}
	return b.sendMessage(userID, chatID, body)
}

// SendTextToUser sends a text message directly to a user.
func (b *Bot) SendTextToUser(userID int64, text string) error {
	return b.SendText(userID, 0, text)
}

// SendTextToChat sends a text message to a chat/group.
func (b *Bot) SendTextToChat(chatID int64, text string) error {
	return b.SendText(0, chatID, text)
}

// SendTextWithKeyboard sends a text message with an inline keyboard.
// Provide either userID (for private message) or chatID (for group message), not both.
func (b *Bot) SendTextWithKeyboard(userID, chatID int64, text string, keyboard *models.InlineKeyboard) error {
	attachments := []models.AttachmentRequest{
		{
			Type: models.AttachmentTypeInlineKeyboard,
			Payload: &models.InlineKeyboardAttachmentRequestPayload{
				Buttons: keyboard.Buttons,
			},
		},
	}

	body := &models.NewMessageBody{
		Text:        &text,
		Attachments: attachments,
	}

	return b.sendMessage(userID, chatID, body)
}

// EditMessage updates the text of an existing message.
func (b *Bot) EditMessage(messageID, newText string) error {
	request := reqMsg.Edit{
		MessageID: messageID,
		NewMessageBody: models.NewMessageBody{
			Text: &newText,
		},
	}

	_, err := b.client.Call(context.Background(), core.EditMsg, request)
	if err != nil && b.Logger != nil {
		b.Logger.Error("bot", "edit_message_failed", fmt.Sprintf("messageID=%s, error=%v", messageID, err))
	}
	return err
}

// DeleteMessage removes a message by its ID.
func (b *Bot) DeleteMessage(messageID string) error {
	request := reqMsg.Delete{
		MessageID: messageID,
	}

	_, err := b.client.Call(context.Background(), core.DeleteMsg, request)
	if err != nil && b.Logger != nil {
		b.Logger.Error("bot", "delete_message_failed", fmt.Sprintf("messageID=%s, error=%v", messageID, err))
	}
	return err
}

func (b *Bot) sendMessage(userID, chatID int64, body *models.NewMessageBody) error {
	request := reqMsg.Send{
		NewMessageBody: *body,
	}

	if chatID != 0 {
		request.ChatID = &chatID
	} else if userID != 0 {
		request.UserID = &userID
	}

	_, err := b.client.Call(context.Background(), core.SendMsg, request)
	if err != nil && b.Logger != nil {
		b.Logger.Error("bot", "send_message_failed", fmt.Sprintf("userID=%d, chatID=%d, error=%v", userID, chatID, err))
	}
	return err
}
