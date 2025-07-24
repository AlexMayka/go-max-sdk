package types

import (
	"context"
	"github.com/AlexMayka/go-max-sdk/types/models"
)

type Transport interface {
	Start(ctx context.Context) (<-chan *models.Update, error)
	Stop() error
}
