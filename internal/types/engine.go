package types

import (
	"github.com/AlexMayka/go-max-sdk/types/models"
)

type BotEngine interface {
	Start() error
	Stop() error
	Dispatch(update *models.Update) error
}
