// Stop bot - demonstrates graceful bot shutdown
package main

import (
	"context"
	maxsdk "github.com/AlexMayka/go-max-sdk"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	bot := maxsdk.NewBot("your-bot-token-here")

	// Configure shutdown timeout
	bot.Config.ShutdownTimeout = 15 * time.Second

	bot.OnCommand("/start", func(ctx *maxsdk.BotContext) {
		ctx.Reply("Stop Bot started! 🟢\n\nCommands:\n/stop - Graceful shutdown\n/ping - Test response")
	})

	bot.OnCommand("/ping", func(ctx *maxsdk.BotContext) {
		ctx.Reply("🏓 Pong! Bot is running...")
	})

	// Manual stop command
	bot.OnCommand("/stop", func(ctx *maxsdk.BotContext) {
		ctx.Reply("🛑 Bot is shutting down gracefully...")

		go func() {
			time.Sleep(100 * time.Millisecond)
			if err := bot.Stop(); err != nil {
				log.Printf("Error stopping bot: %v", err)
			}
		}()
	})

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigChan
		log.Printf("Received signal %v, shutting down...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- bot.Stop()
		}()

		select {
		case err := <-done:
			if err != nil {
				log.Printf("Error during shutdown: %v", err)
			} else {
				log.Println("Bot stopped gracefully")
			}
		case <-ctx.Done():
			log.Println("Shutdown timeout exceeded, forcing exit")
		}

		os.Exit(0)
	}()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Println("Background task: Bot is alive")
			}
		}
	}()

	log.Println("Starting Stop Bot... Press Ctrl+C to stop gracefully")

	if err := bot.Start(); err != nil {
		log.Printf("Bot stopped with error: %v", err)
	} else {
		log.Println("Bot stopped successfully")
	}
}
