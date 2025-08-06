// Config bot - demonstrates bot configuration modification
package main

import (
	"fmt"
	maxsdk "github.com/AlexMayka/go-max-sdk"
	"time"
)

func main() {
	bot := maxsdk.NewBot("your-bot-token-here")

	bot.Config.MaxWorkers = 10
	bot.Config.WorkerQueueSize = 1000
	bot.Config.RequestTimeout = 15 * time.Second
	bot.Config.PollTimeout = 60
	bot.Config.PollLimit = 50
	bot.Config.APITimeout = 10 * time.Second
	bot.Config.APIRateLimit = 30.0
	bot.Config.APIBurstLimit = 100.0

	bot.OnCommand("/start", func(ctx *maxsdk.BotContext) {
		ctx.Reply("Config Bot running with custom settings!")
	})

	bot.OnCommand("/config", func(ctx *maxsdk.BotContext) {
		config := fmt.Sprintf("Current Config:\nMaxWorkers: %d\nWorkerQueueSize: %d\nPollLimit: %d\nAPIRateLimit: %.1f",
			bot.Config.MaxWorkers, bot.Config.WorkerQueueSize, bot.Config.PollLimit, bot.Config.APIRateLimit)
		ctx.Reply(config)
	})

	_ = bot.Start()
}
