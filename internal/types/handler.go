package types

type Handler func(ctx *BotContext)

type Middleware func(handler Handler) Handler