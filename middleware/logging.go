package middleware

import (
	"fmt"
	"github.com/AlexMayka/go-max-sdk/internal/core"
	"time"
)

func Logging(logger core.Logger) core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx *core.BotContext) {
			start := time.Now()
			status := "200 OK"

			defer func() {
				duration := time.Since(start)

				if r := recover(); r != nil {
					status = "500 ERROR"
				}

				timeStr := start.Format("2006/01/02 - 15:04:05")
				username := fmt.Sprintf("user_%d", ctx.UserID)

				message := ctx.Text
				if len(message) > 50 {
					message = message[:47] + "..."
				}
				
				if logger != nil {
					logger.Info("bot", "request",
						fmt.Sprintf("[BOT] %s | %s | %v | %s | %q",
							timeStr, status, duration, username, message))
				}

				if status == "500 ERROR" {
					panic(recover())
				}
			}()

			next(ctx)
		}
	}
}
