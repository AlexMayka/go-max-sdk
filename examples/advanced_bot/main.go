// Advanced bot - demonstrates complex features: custom config, FSM, middleware, routing
package main

import (
	maxsdk "github.com/AlexMayka/go-max-sdk"
	"regexp"
	"strings"
	"time"
)

func main() {
	bot := maxsdk.NewBot("your-bot-token-here")

	// Custom configuration
	bot.Config.MaxWorkers = 15
	bot.Config.WorkerQueueSize = 500
	bot.Config.RequestTimeout = 20 * time.Second
	bot.Config.APIRateLimit = 25.0

	// Create custom router with middleware
	router := maxsdk.DefaultRouter()
	router.Use(loggingMiddleware, rateLimitMiddleware)

	// Public routes
	router.OnCommand("/start", startHandler)
	router.OnCommand("/help", helpHandler)

	// User registration system
	setupRegistrationRoutes(router)

	// Admin panel with authentication
	adminGroup := router.Group("/admin").Use(adminAuthMiddleware)
	setupAdminRoutes(adminGroup)

	// File processing with regex
	router.OnRegex(`^/process\s+(.+)$`, processFileHandler)
	router.OnRegex(`^#(\w+)`, hashtagHandler)

	// Fallback handler
	router.Any(unknownCommandHandler)

	bot.SetRouter(router)
	_ = bot.Start()
}

// === HANDLERS ===

func startHandler(ctx *maxsdk.BotContext) {
	ctx.Reply("🚀 Advanced Bot v2.0\n\nCommands:\n/register - User registration\n/admin/login - Admin access\n/process <text> - Process text\n#hashtag - Track hashtags")
}

func helpHandler(ctx *maxsdk.BotContext) {
	help := "📋 Available Commands:\n\n" +
		"🔹 /start - Start bot\n" +
		"🔹 /register - Begin registration\n" +
		"🔹 /profile - View profile\n" +
		"🔹 /admin/login <pass> - Admin login\n" +
		"🔹 /admin/stats - View stats (admin)\n" +
		"🔹 /process <text> - Process text\n" +
		"🔹 #hashtag - Track hashtag"
	ctx.Reply(help)
}

// === REGISTRATION SYSTEM ===

func setupRegistrationRoutes(router maxsdk.Router) {
	router.OnCommand("/register", registerStartHandler)
	router.UseState("reg_name").Any(registerNameHandler)
	router.UseState("reg_email").Any(registerEmailHandler)
	router.OnCommand("/profile", profileHandler)
}

func registerStartHandler(ctx *maxsdk.BotContext) {
	ctx.Reply("📝 Registration started!\nStep 1/2: Enter your full name:")
	ctx.SetState("reg_name")
}

func registerNameHandler(ctx *maxsdk.BotContext) {
	name := strings.TrimSpace(ctx.Text)
	if len(name) < 2 {
		ctx.Reply("❌ Name too short. Try again:")
		return
	}
	if matched, _ := regexp.MatchString(`^[a-zA-Z\s]+$`, name); !matched {
		ctx.Reply("❌ Name can only contain letters. Try again:")
		return
	}

	ctx.SetData("name", name)
	ctx.Reply("✅ Name: " + name + "\nStep 2/2: Enter your email:")
	ctx.SetState("reg_email")
}

func registerEmailHandler(ctx *maxsdk.BotContext) {
	email := strings.TrimSpace(ctx.Text)
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(emailRegex, email); !matched {
		ctx.Reply("❌ Invalid email format. Try again:")
		return
	}

	ctx.SetData("email", email)
	ctx.SetData("registered", "true")
	name, _ := ctx.GetData("name")
	ctx.Reply("🎉 Registration complete!\n👤 " + name + "\n📧 " + email)
	ctx.SetState("")
}

func profileHandler(ctx *maxsdk.BotContext) {
	if registered, _ := ctx.GetData("registered"); registered != "true" {
		ctx.Reply("❌ Not registered. Use /register first.")
		return
	}
	name, _ := ctx.GetData("name")
	email, _ := ctx.GetData("email")
	ctx.Reply("📋 Your Profile:\n👤 " + name + "\n📧 " + email)
}

// === ADMIN SYSTEM ===

func setupAdminRoutes(adminGroup maxsdk.Router) {
	adminGroup.OnCommand("/stats", adminStatsHandler)
	adminGroup.OnCommand("/users", adminUsersHandler)
}

func adminStatsHandler(ctx *maxsdk.BotContext) {
	ctx.Reply("📊 Admin Stats:\n• Users: 156\n• Messages: 2,341\n• Uptime: 5d 12h")
}

func adminUsersHandler(ctx *maxsdk.BotContext) {
	ctx.Reply("👥 Recent Users:\n• John Doe (active)\n• Jane Smith (offline)\n• Bob Wilson (active)")
}

// === SPECIAL HANDLERS ===

func processFileHandler(ctx *maxsdk.BotContext) {
	text := strings.TrimPrefix(ctx.Text, "/process ")
	processed := strings.ToUpper(text)
	ctx.Reply("🔄 Processed: " + processed)
}

func hashtagHandler(ctx *maxsdk.BotContext) {
	hashtag := ctx.Text[1:] // Remove #
	ctx.Reply("🏷️ Hashtag tracked: #" + hashtag)
}

func unknownCommandHandler(ctx *maxsdk.BotContext) {
	ctx.Reply("❓ Unknown command. Use /help for available commands.")
}

// === MIDDLEWARE ===

func loggingMiddleware(next maxsdk.Handler) maxsdk.Handler {
	return func(ctx *maxsdk.BotContext) {
		next(ctx)
	}
}

func rateLimitMiddleware(next maxsdk.Handler) maxsdk.Handler {
	return func(ctx *maxsdk.BotContext) {
		next(ctx)
	}
}

func adminAuthMiddleware(next maxsdk.Handler) maxsdk.Handler {
	return func(ctx *maxsdk.BotContext) {
		if strings.HasPrefix(ctx.Text, "/admin/login ") {
			password := strings.TrimPrefix(ctx.Text, "/admin/login ")
			if password == "secret123" {
				ctx.SetData("admin", "true")
				ctx.Reply("✅ Admin access granted")
			} else {
				ctx.Reply("❌ Invalid password")
			}
			return
		}

		if admin, _ := ctx.GetData("admin"); admin != "true" {
			ctx.Reply("🔒 Admin access required. Use '/admin/login <password>'")
			return
		}
		next(ctx)
	}
}
