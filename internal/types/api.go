package types

import (
	"context"
	"github.com/AlexMayka/go-max-sdk/internal/config"
)

type APIClient interface {
	Call(ctx context.Context, endpoint config.Endpoint, req interface{}) (interface{}, error)
}