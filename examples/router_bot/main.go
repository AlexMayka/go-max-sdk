// Router bot - demonstrates custom router setup and manual start
package main

import (
	maxsdk "github.com/AlexMayka/go-max-sdk"
)

func main() {
	bot := maxsdk.NewBot("your-bot-token-here")

	router := maxsdk.DefaultRouter()

	router.OnCommand("/start", func(ctx *maxsdk.BotContext) {
		ctx.Reply("Custom router bot started!")
	})
	
	router.OnMessage("hello", func(ctx *maxsdk.BotContext) {
		ctx.Reply("Exact match: hello!")
	})

	router.OnPrefix("!", func(ctx *maxsdk.BotContext) {
		ctx.Reply("Message starts with '!' - " + ctx.Text)
	})

	router.OnSuffix("?", func(ctx *maxsdk.BotContext) {
		ctx.Reply("Question detected: " + ctx.Text)
	})

	router.OnContains("bot", func(ctx *maxsdk.BotContext) {
		ctx.Reply("You mentioned 'bot' in: " + ctx.Text)
	})

	router.OnRegex(`^\d+$`, func(ctx *maxsdk.BotContext) {
		ctx.Reply("You sent a number: " + ctx.Text)
	})

	adminGroup := router.Group("/admin").Use(authMiddleware)
	adminGroup.OnCommand("/status", func(ctx *maxsdk.BotContext) {
		ctx.Reply("Admin status: OK")
	})
	adminGroup.OnCommand("/users", func(ctx *maxsdk.BotContext) {
		ctx.Reply("Admin users count: 42")
	})

	publicGroup := router.Group("/public")
	publicGroup.OnCommand("/help", func(ctx *maxsdk.BotContext) {
		ctx.Reply("Router methods:\n• OnMessage - exact match\n• OnPrefix - starts with\n• OnSuffix - ends with\n• OnContains - contains text\n• OnRegex - regex pattern")
	})
	publicGroup.OnCommand("/info", func(ctx *maxsdk.BotContext) {
		ctx.Reply("Try:\n• 'hello' - exact match\n• '!test' - prefix\n• 'test?' - suffix\n• 'my bot rocks' - contains\n• '123' - regex number")
	})

	bot.SetRouter(router)

	_ = bot.Start()
}

func authMiddleware(next maxsdk.Handler) maxsdk.Handler {
	return func(ctx *maxsdk.BotContext) {
		if adminUser, _ := ctx.GetData("admin"); adminUser != "true" {
			ctx.Reply("Access denied. Use /login first")
			return
		}
		next(ctx)
	}
}
