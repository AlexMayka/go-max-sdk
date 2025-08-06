// Context bot - demonstrates direct message sending, editing, and deletion
package main

import (
	maxsdk "github.com/AlexMayka/go-max-sdk"
	"strconv"
	"strings"
)

func main() {
	bot := maxsdk.NewBot("your-bot-token-here")

	bot.OnCommand("/start", func(ctx *maxsdk.BotContext) {
		ctx.Reply("Context Bot - Direct messaging demo!\n\nCommands:\n/send <user_id> <text> - Send direct message\n/edit <msg_id> <text> - Edit message\n/delete <msg_id> - Delete message\n/broadcast <text> - Send to multiple users")
	})

	bot.OnRegex(`^/send\s+(\d+)\s+(.+)$`, func(ctx *maxsdk.BotContext) {
		parts := strings.SplitN(ctx.Text, " ", 3)
		if len(parts) < 3 {
			ctx.Reply("Usage: /send <user_id> <text>")
			return
		}

		userID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			ctx.Reply("Invalid user ID")
			return
		}

		text := parts[2]
		err = bot.SendTextToUser(userID, text)
		if err != nil {
			ctx.Reply("Failed to send message: " + err.Error())
		} else {
			ctx.Reply("Message sent to user " + parts[1])
		}
	})

	bot.OnRegex(`^/edit\s+(\S+)\s+(.+)$`, func(ctx *maxsdk.BotContext) {
		parts := strings.SplitN(ctx.Text, " ", 3)
		if len(parts) < 3 {
			ctx.Reply("Usage: /edit <message_id> <new_text>")
			return
		}

		messageID := parts[1]
		newText := parts[2]

		err := ctx.EditMessage(messageID, newText)
		if err != nil {
			ctx.Reply("Failed to edit message: " + err.Error())
		} else {
			ctx.Reply("Message " + messageID + " edited!")
		}
	})

	bot.OnRegex(`^/delete\s+(\S+)$`, func(ctx *maxsdk.BotContext) {
		parts := strings.Fields(ctx.Text)
		if len(parts) < 2 {
			ctx.Reply("Usage: /delete <message_id>")
			return
		}

		messageID := parts[1]
		err := ctx.DeleteMessage(messageID)
		if err != nil {
			ctx.Reply("Failed to delete message: " + err.Error())
		} else {
			ctx.Reply("Message " + messageID + " deleted!")
		}
	})

	bot.OnRegex(`^/broadcast\s+(.+)$`, func(ctx *maxsdk.BotContext) {
		text := strings.TrimPrefix(ctx.Text, "/broadcast ")

		userIDs := []int64{12345, 67890, 11111}

		sent := 0
		for _, userID := range userIDs {
			err := bot.SendTextToUser(userID, "📢 Broadcast: "+text)
			if err == nil {
				sent++
			}
		}

		ctx.Reply("Broadcast sent to " + strconv.Itoa(sent) + " users")
	})

	bot.OnCommand("/info", func(ctx *maxsdk.BotContext) {
		info := "Message Context Info:\n" +
			"👤 User ID: " + strconv.FormatInt(ctx.UserID, 10) + "\n" +
			"💬 Chat ID: " + strconv.FormatInt(ctx.ChatID, 10) + "\n" +
			"📝 Text: " + ctx.Text + "\n" +
			"👋 First Name: " + ctx.FirstName

		if ctx.LastName != "" {
			info += "\n🔖 Last Name: " + ctx.LastName
		}
		if ctx.Username != "" {
			info += "\n📛 Username: @" + ctx.Username
		}

		ctx.Reply(info)
	})

	bot.OnCommand("/demo", func(ctx *maxsdk.BotContext) {
		_ = bot.SendTextToUser(ctx.UserID, "This message was sent using bot.SendTextToUser()")

		if ctx.ChatID != 0 {
			_ = bot.SendTextToChat(ctx.ChatID, "This message was sent using bot.SendTextToChat()")
		}

		ctx.Reply("Demo messages sent!")
	})

	_ = bot.Start()
}
