package engine

import (
	"context"
	"sync"
	"time"

	"github.com/AlexMayka/go-max-sdk/internal/types"
	"github.com/AlexMayka/go-max-sdk/types/models"
	"golang.org/x/time/rate"
)

type botEngine struct {
	registry  types.RouteRegistry
	fsm       types.FSM
	client    types.APIClient
	transport types.Transport
	ctx       context.Context
	cancel    context.CancelFunc

	workerPool  chan *models.Update
	cnf         types.EngineConfig
	wg          sync.WaitGroup
	rateLimiter *rate.Limiter
}

func NewBotEngine(registry types.RouteRegistry, fsm types.FSM, client types.APIClient, transport types.Transport, ctx context.Context, cnf types.EngineConfig) types.BotEngine {
	ctx, cancel := context.WithCancel(ctx)

	var limiter *rate.Limiter
	if cnf.RateLimit > 0 {
		limiter = rate.NewLimiter(rate.Every(cnf.RatePeriod/time.Duration(cnf.RateLimit)), cnf.BurstLimit)
	}

	return &botEngine{
		registry:    registry,
		fsm:         fsm,
		client:      client,
		transport:   transport,
		ctx:         ctx,
		cancel:      cancel,
		cnf:         cnf,
		wg:          sync.WaitGroup{},
		rateLimiter: limiter,
	}
}

func (e *botEngine) worker(id int, updates <-chan *models.Update) {
	defer e.wg.Done()

	for update := range updates {
		handler := e.processUpdate(update)
		if handler == nil {
			continue
		}

		e.executeHandlerSimple(handler, update)
	}
}

func (e *botEngine) Start() error {
	e.workerPool = make(chan *models.Update, e.cnf.WorkerQueueSize)

	ch, err := e.transport.Start(e.ctx)
	if err != nil {
		return err
	}

	for i := 0; i < e.cnf.MaxWorkers; i++ {
		e.wg.Add(1)
		go e.worker(i, e.workerPool)
	}

	e.wg.Add(1)
	go e.messageReader(ch)

	return nil
}

func (e *botEngine) messageReader(ch <-chan *models.Update) {
	defer e.wg.Done()
	defer close(e.workerPool)

	for {
		select {
		case update, ok := <-ch:
			if !ok {
				return
			}
			if update == nil {
				continue
			}

			e.dispatch(update)

		case <-e.ctx.Done():
			return
		}
	}
}

func (e *botEngine) dispatch(update *models.Update) {
	if e.rateLimiter != nil && !e.rateLimiter.Allow() {
		return
	}

	select {
	case e.workerPool <- update:
	case <-e.ctx.Done():
	default:
	}
}

func (e *botEngine) Stop() error {
	e.cancel()

	done := make(chan struct{})
	go func() {
		e.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(e.cnf.ShutdownTimeout):
	}

	return e.transport.Stop()
}
