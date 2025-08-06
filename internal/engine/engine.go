// Package engine provides the core bot execution engine with worker pool architecture.
// It manages message routing, worker coordination, and update processing for the bot.
package engine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AlexMayka/go-max-sdk/internal/core"
	"github.com/AlexMayka/go-max-sdk/types/models"
)

// botEngine implements the core.BotEngine interface and manages the bot's execution lifecycle.
// It coordinates between transport layer, route registry, FSM, and worker pool for processing updates.
type botEngine struct {
	registry  core.RouteRegistry
	fsm       core.FSM
	client    core.APIClient
	transport core.Transport
	ctx       context.Context
	cancel    context.CancelFunc

	workerPool chan *models.Update
	wg         sync.WaitGroup

	maxWorkers      int
	workerQueueSize int
	requestTimeout  time.Duration
	shutdownTimeout time.Duration
	retryAttempts   int
	retryDelay      time.Duration

	logger core.Logger
}

// NewBotEngine creates a new bot engine instance with the specified configuration.
// The engine uses a worker pool architecture to process incoming updates concurrently.
// Parameters:
//   - registry: route registry containing all registered handlers
//   - fsm: finite state machine for managing user conversation states
//   - client: API client for making requests to the bot API
//   - transport: transport layer (long polling, webhook, etc.)
//   - ctx: parent context for cancellation
//   - maxWorkers: number of worker goroutines to spawn
//   - workerQueueSize: buffer size for the worker pool channel
//   - requestTimeout: timeout for individual handler execution
//   - shutdownTimeout: timeout for graceful shutdown
//   - retryAttempts: number of retry attempts for failed handlers
//   - retryDelay: delay between retry attempts
//   - logger: logger instance for structured logging
func NewBotEngine(registry core.RouteRegistry, fsm core.FSM, client core.APIClient, transport core.Transport, ctx context.Context, maxWorkers int, workerQueueSize int, requestTimeout time.Duration, shutdownTimeout time.Duration, retryAttempts int, retryDelay time.Duration, logger core.Logger) core.BotEngine {
	ctx, cancel := context.WithCancel(ctx)

	return &botEngine{
		registry:        registry,
		fsm:             fsm,
		client:          client,
		transport:       transport,
		ctx:             ctx,
		cancel:          cancel,
		wg:              sync.WaitGroup{},
		maxWorkers:      maxWorkers,
		workerQueueSize: workerQueueSize,
		requestTimeout:  requestTimeout,
		shutdownTimeout: shutdownTimeout,
		retryAttempts:   retryAttempts,
		retryDelay:      retryDelay,
		logger:          logger,
	}
}

// worker processes updates from the worker pool channel.
// Each worker runs in its own goroutine and processes updates sequentially.
func (e *botEngine) worker(updates <-chan *models.Update) {
	defer e.wg.Done()

	for update := range updates {
		handler := e.processUpdate(update)
		if handler == nil {
			continue
		}

		e.executeHandlerSimple(handler, update)
	}
}

// Start initializes and starts the bot engine.
// It creates the worker pool, starts the transport layer, spawns workers, and begins message processing.
func (e *botEngine) Start() error {
	e.workerPool = make(chan *models.Update, e.workerQueueSize)

	if e.logger != nil {
		e.logger.Info("engine", "starting", fmt.Sprintf("workers=%d, queue_size=%d", e.maxWorkers, e.workerQueueSize))
	}

	ch, err := e.transport.Start(e.ctx)
	if err != nil {
		if e.logger != nil {
			e.logger.Error("engine", "transport_start_failed", fmt.Sprintf("error=%v", err))
		}
		return err
	}

	for i := 0; i < e.maxWorkers; i++ {
		e.wg.Add(1)
		go e.worker(e.workerPool)
	}

	e.wg.Add(1)
	go e.messageReader(ch)

	if e.logger != nil {
		e.logger.Info("engine", "started", "bot_engine_ready")
	}

	return nil
}

// messageReader reads updates from the transport layer and dispatches them to workers.
// It runs in a separate goroutine and handles graceful shutdown.
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

// dispatch sends an update to the worker pool for processing.
// Uses non-blocking send to avoid blocking the message reader.
func (e *botEngine) dispatch(update *models.Update) {
	select {
	case e.workerPool <- update:
	case <-e.ctx.Done():
	default:
	}
}

// Stop gracefully shuts down the bot engine.
// It cancels the context, waits for workers to finish (with timeout), and stops the transport.
func (e *botEngine) Stop() error {
	if e.logger != nil {
		e.logger.Info("engine", "stopping", "graceful_shutdown_initiated")
	}

	e.cancel()

	done := make(chan struct{})
	go func() {
		e.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		if e.logger != nil {
			e.logger.Info("engine", "workers_stopped", "all_workers_finished")
		}
	case <-time.After(e.shutdownTimeout):
		if e.logger != nil {
			e.logger.Warn("engine", "shutdown_timeout", fmt.Sprintf("timeout=%v", e.shutdownTimeout))
		}
	}

	err := e.transport.Stop()
	if err != nil && e.logger != nil {
		e.logger.Error("engine", "transport_stop_failed", fmt.Sprintf("error=%v", err))
	}

	if e.logger != nil {
		e.logger.Info("engine", "stopped", "bot_engine_shutdown_complete")
	}

	return err
}
