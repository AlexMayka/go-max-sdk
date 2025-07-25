package types

type BotEngine interface {
	Start() error
	Stop() error
}
