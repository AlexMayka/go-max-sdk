package engine

import (
	"context"
	"testing"

	"github.com/AlexMayka/go-max-sdk/internal/fsm/local"
	"github.com/AlexMayka/go-max-sdk/internal/registry"
	"github.com/AlexMayka/go-max-sdk/internal/router"
	"github.com/AlexMayka/go-max-sdk/internal/types"
	"github.com/AlexMayka/go-max-sdk/types/models"
)

func BenchmarkDispatcher(b *testing.B) {
	mockHandler := func(ctx *types.BotContext) {}

	update := &models.Update{
		UpdateType: models.UpdateTypeMessageCreated,
		Message: &models.Message{
			Sender: &models.User{UserID: 123},
			Body:   &models.MessageBody{Text: "hello"},
		},
	}

	b.Run("ProcessUpdate_NoState", func(b *testing.B) {
		root := router.NewRouter("", nil)
		root.OnMessage("hello", mockHandler)
		root.OnMessage("hi", mockHandler)
		root.OnMessage("hey", mockHandler)

		reg := registry.BuildFrom(root)
		fsm := local.NewFSM()

		engine := &botEngine{
			registry: reg,
			fsm:      fsm,
		}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			handler := engine.processUpdate(update)
			_ = handler
		}
	})

	b.Run("ProcessUpdate_WithState", func(b *testing.B) {
		root := router.NewRouter("", nil)
		root.UseState("waiting").OnMessage("hello", mockHandler)
		root.UseState("waiting").OnMessage("response", mockHandler)

		reg := registry.BuildFrom(root)
		fsm := local.NewFSM()
		fsm.SetStateUser(123, "waiting")

		engine := &botEngine{
			registry: reg,
			fsm:      fsm,
		}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			handler := engine.processUpdate(update)
			_ = handler
		}
	})

	b.Run("ProcessUpdate_CommandDetection", func(b *testing.B) {
		commandUpdate := &models.Update{
			UpdateType: models.UpdateTypeMessageCreated,
			Message: &models.Message{
				Sender: &models.User{UserID: 123},
				Body:   &models.MessageBody{Text: "/start"},
			},
		}

		root := router.NewRouter("", nil)
		root.OnCommand("/start", mockHandler)
		root.OnCommand("/help", mockHandler)
		root.OnCommand("/settings", mockHandler)

		reg := registry.BuildFrom(root)
		fsm := local.NewFSM()

		engine := &botEngine{
			registry: reg,
			fsm:      fsm,
		}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			handler := engine.processUpdate(commandUpdate)
			_ = handler
		}
	})

	b.Run("ProcessUpdate_FallbackToAny", func(b *testing.B) {
		unmatchUpdate := &models.Update{
			UpdateType: models.UpdateTypeMessageCreated,
			Message: &models.Message{
				Sender: &models.User{UserID: 123},
				Body:   &models.MessageBody{Text: "random message"},
			},
		}

		root := router.NewRouter("", nil)
		root.OnMessage("hello", mockHandler)
		root.Any(mockHandler)

		reg := registry.BuildFrom(root)
		fsm := local.NewFSM()

		engine := &botEngine{
			registry: reg,
			fsm:      fsm,
		}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			handler := engine.processUpdate(unmatchUpdate)
			_ = handler
		}
	})

	b.Run("GetEventType", func(b *testing.B) {
		engine := &botEngine{}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			eventType := engine.getEventType(update)
			_ = eventType
		}
	})

	b.Run("GetEventType_Command", func(b *testing.B) {
		commandUpdate := &models.Update{
			UpdateType: models.UpdateTypeMessageCreated,
			Message: &models.Message{
				Body: &models.MessageBody{Text: "/start"},
			},
		}

		engine := &botEngine{}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			eventType := engine.getEventType(commandUpdate)
			_ = eventType
		}
	})
}

func BenchmarkMatching(b *testing.B) {
	pattern := "hello"
	message := "hello world"
	handler := newMockRouteHandler(types.MatchPrefix, pattern, nil, types.EventMessage)

	b.Run("CheckPrefix", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			result := checkPrefix(&message, handler)
			_ = result
		}
	})

	b.Run("CheckExact", func(b *testing.B) {
		exactMsg := "hello"
		exactHandler := newMockRouteHandler(types.MatchExact, pattern, nil, types.EventMessage)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			result := checkExact(&exactMsg, exactHandler)
			_ = result
		}
	})

	b.Run("CheckContains", func(b *testing.B) {
		containsMsg := "this is a hello message"
		containsHandler := newMockRouteHandler(types.MatchContains, pattern, nil, types.EventMessage)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			result := checkContains(&containsMsg, containsHandler)
			_ = result
		}
	})

	b.Run("CheckRegex", func(b *testing.B) {
		regexPattern := `^\d+$`
		regexMsg := "12345"
		regexHandler := newMockRouteHandler(types.MatchRegex, regexPattern, nil, types.EventMessage)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			result := checkRegex(&regexMsg, regexHandler)
			_ = result
		}
	})
}

func BenchmarkFullPipeline(b *testing.B) {
	b.Run("CompleteMessageProcessing", func(b *testing.B) {
		root := router.NewRouter("", nil)

		root.OnMessage("hello", func(ctx *types.BotContext) {})
		root.OnMessage("hi", func(ctx *types.BotContext) {})
		root.OnPrefix("start", func(ctx *types.BotContext) {})
		root.OnRegex(`^\d+$`, func(ctx *types.BotContext) {})

		root.OnCommand("/help", func(ctx *types.BotContext) {})
		root.OnCommand("/settings", func(ctx *types.BotContext) {})

		root.UseState("waiting").OnMessage("response", func(ctx *types.BotContext) {})
		root.UseState("menu").OnCallback("button_1", func(ctx *types.BotContext) {})

		root.Any(func(ctx *types.BotContext) {})

		reg := registry.BuildFrom(root)
		fsm := local.NewFSM()
		cnf := types.DefaultEngineConfig()

		engine := NewBotEngine(
			reg, fsm, nil, nil,
			context.Background(), *cnf,
		).(*botEngine)

		update := &models.Update{
			UpdateType: models.UpdateTypeMessageCreated,
			Message: &models.Message{
				Sender: &models.User{UserID: 123},
				Body:   &models.MessageBody{Text: "hello"},
			},
		}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			handler := engine.processUpdate(update)
			_ = handler
		}
	})

	b.Run("LargeRouterDispatch", func(b *testing.B) {
		root := router.NewRouter("", nil)

		for i := 0; i < 1000; i++ {
			root.OnMessage("message", func(ctx *types.BotContext) {})
		}

		reg := registry.BuildFrom(root)
		fsm := local.NewFSM()

		engine := &botEngine{
			registry: reg,
			fsm:      fsm,
		}

		update := &models.Update{
			UpdateType: models.UpdateTypeMessageCreated,
			Message: &models.Message{
				Sender: &models.User{UserID: 123},
				Body:   &models.MessageBody{Text: "message"},
			},
		}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			handler := engine.processUpdate(update)
			_ = handler
		}
	})
}

func BenchmarkMemoryUsage(b *testing.B) {
	b.Run("UpdateObjectCreation", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			update := &models.Update{
				UpdateType: models.UpdateTypeMessageCreated,
				Message: &models.Message{
					Sender: &models.User{UserID: 123},
					Body:   &models.MessageBody{Text: "hello"},
				},
			}
			_ = update
		}
	})

	b.Run("HandlerChainExecution", func(b *testing.B) {
		root := router.NewRouter("", nil)

		middleware1 := func(next types.Handler) types.Handler {
			return func(ctx *types.BotContext) { next(ctx) }
		}
		middleware2 := func(next types.Handler) types.Handler {
			return func(ctx *types.BotContext) { next(ctx) }
		}

		root.Use(middleware1, middleware2).OnMessage("test", func(ctx *types.BotContext) {})

		reg := registry.BuildFrom(root)
		fsm := local.NewFSM()
		cnf := types.DefaultEngineConfig()

		engine := NewBotEngine(
			reg, fsm, nil, nil,
			context.Background(), *cnf,
		).(*botEngine)

		update := &models.Update{
			UpdateType: models.UpdateTypeMessageCreated,
			Message: &models.Message{
				Sender: &models.User{UserID: 123},
				Body:   &models.MessageBody{Text: "test"},
			},
		}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			handler := engine.processUpdate(update)
			if handler != nil {
				ctx := &types.BotContext{
					Ctx:    context.Background(),
					UserID: 123,
					Text:   "test",
				}
				handler(ctx)
			}
		}
	})
}
