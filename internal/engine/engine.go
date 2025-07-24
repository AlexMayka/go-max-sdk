package engine

import (
	"context"
	"github.com/AlexMayka/go-max-sdk/internal/registry"
	"github.com/AlexMayka/go-max-sdk/internal/types"
	"github.com/AlexMayka/go-max-sdk/types/models"
)

type botEngine struct {
	registry  types.RouteRegistry
	fsm       types.FSM
	client    types.APIClient
	transport types.Transport
	ctx       context.Context
}

func NewBotEngine(router types.Router, fsm types.FSM, client types.APIClient, transport types.Transport, ctx context.Context) types.BotEngine {
	return &botEngine{
		registry:  registry.BuildFrom(router),
		fsm:       fsm,
		client:    client,
		transport: transport,
		ctx:       ctx,
	}
}

func (e *botEngine) Start() error {
	return nil
}

func (e *botEngine) Stop() error {
	return nil
}

func (e *botEngine) Dispatch(update *models.Update) error {
	return nil
}
