// Package transport provides communication layer implementations for receiving bot updates.
package transport

import (
	"context"
	"fmt"
	"github.com/AlexMayka/go-max-sdk/internal/core"
	"time"

	"github.com/AlexMayka/go-max-sdk/types/models"
	req "github.com/AlexMayka/go-max-sdk/types/requests/subscriptions"
	res "github.com/AlexMayka/go-max-sdk/types/responses/subscriptions"
)

// LongPolling implements the core.Transport interface using HTTP long polling.
// It continuously polls the API for new updates and delivers them through a channel.
// Features include:
//   - Automatic retry logic with exponential backoff
//   - Context-based cancellation
//   - Update marker tracking for reliable delivery
//   - Configurable polling parameters
type LongPolling struct {
	ctx    context.Context
	cancel context.CancelFunc

	client core.APIClient

	pollTimeout    int
	pollLimit      int
	updatesBuffer  int
	pollRetryDelay time.Duration
	pollMaxRetries int
	marker         int64

	logger core.Logger
}

// NewLongPolling creates a new long polling transport with the specified configuration.
// Parameters:
//   - pollTimeout: server-side timeout for long polling requests (seconds)
//   - pollLimit: maximum number of updates to fetch in one request
//   - updatesBuffer: buffer size for the update channel
//   - pollRetryDelay: base delay for retry attempts
//   - pollMaxRetries: maximum retry attempts before capping delay
//   - client: API client for making requests
//   - logger: logger instance for structured logging
func NewLongPolling(pollTimeout int, pollLimit int, updatesBuffer int, pollRetryDelay time.Duration, pollMaxRetries int, client core.APIClient, logger core.Logger) core.Transport {
	return &LongPolling{
		pollTimeout:    pollTimeout,
		pollLimit:      pollLimit,
		updatesBuffer:  updatesBuffer,
		pollRetryDelay: pollRetryDelay,
		pollMaxRetries: pollMaxRetries,
		client:         client,
		logger:         logger,
	}
}

// Start begins the long polling process and returns a channel for receiving updates.
// The polling runs in a separate goroutine and can be stopped via context cancellation or Stop().
// It implements automatic retry logic with exponential backoff for resilient operation.
func (l *LongPolling) Start(ctx context.Context) (<-chan *models.Update, error) {
	if l.ctx != nil && l.ctx.Err() == nil {
		return nil, fmt.Errorf("transport already started")
	}

	l.ctx, l.cancel = context.WithCancel(ctx)
	ch := make(chan *models.Update, l.updatesBuffer)

	if l.logger != nil {
		l.logger.Info("transport", "polling_started",
			fmt.Sprintf("timeout=%d, limit=%d, buffer=%d", l.pollTimeout, l.pollLimit, l.updatesBuffer))
	}

	go func() {
		defer close(ch)
		defer l.cancel()

		typesUpdate := convertUpdateType(models.UpdateTypeMessageCreated, models.UpdateTypeMessageCallback, models.UpdateBotStarted)
		reqModel := &req.GetUpdates{
			Limit:   &l.pollLimit,
			Timeout: &l.pollTimeout,
			Marker:  &l.marker,
			Types:   typesUpdate,
		}

		retryCount := 0

		for {
			select {
			case <-l.ctx.Done():
				return
			default:
				response, err := l.client.Call(l.ctx, core.GetSubscribeUpdate, reqModel)
				if err != nil {
					retryCount++

					if l.logger != nil {
						l.logger.Warn("transport", "polling_retry",
							fmt.Sprintf("error=%v, retry=%d/%d", err, retryCount, l.pollMaxRetries))
					}

					if retryCount > l.pollMaxRetries {
						if l.logger != nil {
							l.logger.Error("transport", "polling_failed",
								fmt.Sprintf("exceeded %d retries, last error: %v", l.pollMaxRetries, err))
						}

						retryCount = l.pollMaxRetries
					}

					delay := l.pollRetryDelay * time.Duration(retryCount)
					if delay > 30*time.Second {
						delay = 30 * time.Second
					}

					select {
					case <-time.After(delay):
						continue
					case <-l.ctx.Done():
						if l.logger != nil {
							l.logger.Debug("transport", "polling_cancelled", "context cancelled during retry")
						}
						return
					}
				}

				retryCount = 0

				getUpdates, ok := response.(res.GetUpdates)
				if !ok {
					if l.logger != nil {
						l.logger.Error("transport", "invalid_response_type",
							fmt.Sprintf("expected res.GetUpdates, got %T", response))
					}
					continue
				}

				if l.logger != nil && len(getUpdates.Updates) > 0 {
					l.logger.Debug("transport", "updates_received",
						fmt.Sprintf("count=%d", len(getUpdates.Updates)))
				}
				for _, update := range getUpdates.Updates {
					select {
					case ch <- &update:
					case <-l.ctx.Done():
						return
					}
				}

				if getUpdates.Marker != nil {
					l.marker = *getUpdates.Marker
					reqModel.Marker = &l.marker
				}
			}
		}

	}()

	return ch, nil
}

// Stop gracefully stops the long polling transport by cancelling the context.
func (l *LongPolling) Stop() error {
	if l.cancel != nil {
		l.cancel()
		if l.logger != nil {
			l.logger.Info("transport", "polling_stopped", "")
		}
	}
	return nil
}

// convertUpdateType converts model update types to string array for API requests.
func convertUpdateType(ups ...models.UpdateType) []string {
	result := make([]string, 0, len(ups))
	for _, up := range ups {
		result = append(result, string(up))
	}
	return result
}
