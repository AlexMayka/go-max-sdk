// Echo bot - repeats user messages
package main

import (
	maxsdk "github.com/AlexMayka/go-max-sdk"
)

func main() {
	bot := maxsdk.NewBot("your-bot-token-here")

	bot.Any(func(ctx *maxsdk.BotContext) {
		if ctx.Text != "" {
			ctx.Reply("You said: " + ctx.Text)
		}
	})

	_ = bot.Start()
}
