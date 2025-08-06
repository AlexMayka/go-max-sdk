// Logger bot - demonstrates custom logger (simplified for public API)
package main

import (
	"fmt"
	maxsdk "github.com/AlexMayka/go-max-sdk"
	"log"
)

func main() {
	bot := maxsdk.NewBot("your-bot-token-here")

	bot.OnCommand("/start", func(ctx *maxsdk.BotContext) {
		bot.Logger.Info("bot", "start_command", "User started the bot")
		ctx.Reply("Logger bot started! Check console for logs.")
	})

	bot.OnCommand("/log", func(ctx *maxsdk.BotContext) {
		bot.Logger.Debug("bot", "debug_test", "Debug message from user")
		bot.Logger.Info("bot", "info_test", "Info message from user")
		bot.Logger.Warn("bot", "warn_test", "Warning message from user")
		bot.Logger.Error("bot", "error_test", "Error message from user")
		ctx.Reply("Logged messages at all levels! Check console.")
	})

	bot.OnCommand("/stats", func(ctx *maxsdk.BotContext) {
		userInfo := fmt.Sprintf("Stats requested by user %d (%s)", ctx.UserID, ctx.FirstName)
		bot.Logger.Info("bot", "stats_requested", userInfo)

		ctx.Reply("📊 Bot Stats:\n• Active since startup\n• Logger: Console\n• Debug info written to console")
	})

	bot.OnCommand("/disable", func(ctx *maxsdk.BotContext) {
		bot.Logging = false
		ctx.Reply("🔇 Bot logging disabled")
	})

	bot.OnCommand("/enable", func(ctx *maxsdk.BotContext) {
		bot.Logging = true
		bot.Logger.Info("bot", "logging_enabled", "User re-enabled logging")
		ctx.Reply("🔊 Bot logging enabled")
	})

	log.Println("Starting logger bot...")
	_ = bot.Start()
}
