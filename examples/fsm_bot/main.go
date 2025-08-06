// FSM bot - demonstrates finite state machine with separate handler functions
package main

import (
	maxsdk "github.com/AlexMayka/go-max-sdk"
)

func main() {
	bot := maxsdk.NewBot("your-bot-token-here")

	bot.OnCommand("/start", startHandler)
	bot.OnCommand("/order", orderHandler)
	bot.UseState("choosing_pizza").Any(pizzaChoiceHandler)
	bot.UseState("choosing_size").Any(sizeChoiceHandler)
	bot.UseState("awaiting_address").Any(addressHandler)

	_ = bot.Start()
}

func startHandler(ctx *maxsdk.BotContext) {
	ctx.Reply("Welcome to Pizza Bot! 🍕\nUse /order to start ordering")
}

func orderHandler(ctx *maxsdk.BotContext) {
	ctx.Reply("Choose your pizza:\n1. Margherita\n2. Pepperoni\n3. Hawaiian")
	ctx.SetState("choosing_pizza")
}

func pizzaChoiceHandler(ctx *maxsdk.BotContext) {
	var pizza string
	switch ctx.Text {
	case "1":
		pizza = "Margherita"
	case "2":
		pizza = "Pepperoni"
	case "3":
		pizza = "Hawaiian"
	default:
		ctx.Reply("Please choose 1, 2, or 3")
		return
	}
	
	ctx.SetData("pizza", pizza)
	ctx.Reply("Great! You chose " + pizza + "\nNow choose size:\n1. Small ($10)\n2. Medium ($15)\n3. Large ($20)")
	ctx.SetState("choosing_size")
}

func sizeChoiceHandler(ctx *maxsdk.BotContext) {
	var size, price string
	switch ctx.Text {
	case "1":
		size, price = "Small", "$10"
	case "2":
		size, price = "Medium", "$15"
	case "3":
		size, price = "Large", "$20"
	default:
		ctx.Reply("Please choose 1, 2, or 3")
		return
	}
	
	ctx.SetData("size", size)
	ctx.SetData("price", price)
	pizza, _ := ctx.GetData("pizza")
	ctx.Reply("Perfect! " + size + " " + pizza + " for " + price + "\nPlease enter your delivery address:")
	ctx.SetState("awaiting_address")
}

func addressHandler(ctx *maxsdk.BotContext) {
	if len(ctx.Text) < 10 {
		ctx.Reply("Please enter a complete address (at least 10 characters)")
		return
	}
	
	pizza, _ := ctx.GetData("pizza")
	size, _ := ctx.GetData("size")
	price, _ := ctx.GetData("price")
	
	ctx.Reply("Order confirmed! 🎉\n" +
		"Pizza: " + pizza + "\n" +
		"Size: " + size + "\n" +
		"Price: " + price + "\n" +
		"Address: " + ctx.Text + "\n\n" +
		"Delivery time: 30 minutes\nUse /order to order again!")
	
	ctx.SetState("")
}