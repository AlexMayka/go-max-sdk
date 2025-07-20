package maxsdk

import (
	"github.com/AlexMayka/go-max-sdk/fsm"
	"github.com/AlexMayka/go-max-sdk/logger"
	"github.com/AlexMayka/go-max-sdk/transport"
)

type Option func(bot *Bot)

func WithFMS(fms fsm.FMS) Option {
	return func(bot *Bot) { bot.fsm = fms }
}

func WithTransport(transport transport.Transport) Option {
	return func(bot *Bot) { bot.transport = transport }
}

func WithLogger(logger logger.Logger) Option {
	return func(bot *Bot) { bot.logger = logger }
}
