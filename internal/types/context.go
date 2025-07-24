package types

import "context"

type BotContext struct {
	ctx context.Context

	UserID int64
	ChatID int64
}
