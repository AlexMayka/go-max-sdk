package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexMayka/go-max-sdk/internal/core"
	"io"
	"math"
	"net/http"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

// Client implements the core.APIClient interface for making HTTP requests to the MAX Bot API.
// It handles authentication, rate limiting, retries, and request/response marshaling.
type Client struct {
	token      string
	logger     core.Logger
	httpClient *http.Client

	// Token bucket rate limiter fields
	burstLimit float64    // Maximum tokens in bucket
	rateLimit  float64    // Tokens per second refill rate
	tokens     uint64     // Current tokens available (stored as uint64 for atomic)
	lastRefill int64      // Last refill time in nanoseconds (for atomic)
	refillMu   sync.Mutex // Mutex only for refill operations

	// Retry configuration
	maxRetries int           // Maximum retry attempts
	retryDelay time.Duration // Base delay between retries
}

// NewClient creates a new API client instance with the specified configuration.
// Parameters:
//   - token: Bot authentication token
//   - logger: Logger instance for request/response logging
//   - apiTimeout: HTTP request timeout duration
//   - burstLimit: Maximum burst size for token bucket rate limiting
//   - rateLimit: Token refill rate per second (0 = no limit)
//   - maxRetries: Maximum retry attempts for failed requests
//   - retryDelay: Base delay between retry attempts
func NewClient(token string, logger core.Logger, apiTimeout time.Duration, burstLimit float64, rateLimit float64, maxRetries int, retryDelay time.Duration) core.APIClient {
	now := time.Now()
	return &Client{
		token:      token,
		logger:     logger,
		burstLimit: burstLimit,
		rateLimit:  rateLimit,
		tokens:     math.Float64bits(burstLimit),
		lastRefill: now.UnixNano(),
		maxRetries: maxRetries,
		retryDelay: retryDelay,
		httpClient: &http.Client{
			Timeout: apiTimeout,
		},
	}
}

func (c *Client) sendRequest(ctx context.Context, jsonBody interface{}, cfg *EndpointConfig, addr string) (interface{}, error) {
	var body []byte
	var err error

	if jsonBody != nil {
		body, err = json.Marshal(jsonBody)
		if err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequestWithContext(ctx, string(cfg.Method), addr, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if len(body) > 0 {
		request.Header.Set(string(Context), string(cfg.ContentType))
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := response.Body.Close(); err != nil && c.logger != nil {
			c.logger.Warn("api_client", "body_close_error", fmt.Sprintf("error=%v", err))
		}
	}()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if c.logger != nil {
			c.logger.Error("api_client", "http_error", fmt.Sprintf("status=%d, body=%s", response.StatusCode, string(responseBody)))
		}
		return nil, fmt.Errorf("unexpected status code: %d - %s", response.StatusCode, string(responseBody))
	}

	responsePtr := reflect.New(cfg.ResponseModel)
	responseValue := responsePtr.Interface()

	err = json.Unmarshal(responseBody, responseValue)
	if err != nil {
		return nil, err
	}

	return responseValue, nil
}

// Call executes an API request with the given endpoint and request data.
// It handles request parsing, URL building, rate limiting, retries, and response unmarshaling.
// Returns the response data or an error if the request fails after all retry attempts.
func (c *Client) Call(ctx context.Context, endpoint core.Endpoint, req interface{}) (interface{}, error) {
	cfg := EndpointConfigs[endpoint]

	pathParams, queryParams, jsonBody, err := c.parseRequest(req, cfg)
	if err != nil {
		return nil, err
	}

	addr := c.buildURL(cfg.Path, pathParams, queryParams)

	c.waitForRateLimit()

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(attempt) * c.retryDelay
			if delay > 5*time.Second {
				delay = 5 * time.Second
			}

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		if c.logger != nil && attempt == 0 {
			c.logger.Debug("api_client", "request_start", fmt.Sprintf("endpoint=%s, url=%s", endpoint, addr))
		}

		answer, err := c.sendRequest(ctx, jsonBody, cfg, addr)
		if err == nil {
			if c.logger != nil {
				c.logger.Debug("api_client", "request_success", fmt.Sprintf("endpoint=%s, attempt=%d", endpoint, attempt+1))
			}
			return answer, nil
		}

		if c.logger != nil {
			if attempt < c.maxRetries {
				c.logger.Warn("api_client", "retry_attempt", fmt.Sprintf("endpoint=%s, attempt=%d/%d, error=%v", endpoint, attempt+1, c.maxRetries+1, err))
			} else {
				c.logger.Error("api_client", "request_failed", fmt.Sprintf("endpoint=%s, exhausted_retries=%d, error=%v", endpoint, c.maxRetries+1, err))
			}
		}

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts", c.maxRetries+1)
}

// waitForRateLimit implements a lock-free token bucket rate limiting mechanism.
// It uses atomic operations for performance and only locks during token refill.
func (c *Client) waitForRateLimit() {
	if c.rateLimit <= 0 || c.burstLimit <= 0 {
		return
	}

	for {
		currentBits := atomic.LoadUint64(&c.tokens)
		currentTokens := math.Float64frombits(currentBits)

		if currentTokens >= 1.0 {
			newTokens := currentTokens - 1.0
			newBits := math.Float64bits(newTokens)

			if atomic.CompareAndSwapUint64(&c.tokens, currentBits, newBits) {
				if c.logger != nil {
					c.logger.Debug("api_client", "rate_limit_fast_path", fmt.Sprintf("tokens_remaining=%.2f", newTokens))
				}
				return
			}
			continue
		}

		c.refillTokens()

		currentBits = atomic.LoadUint64(&c.tokens)
		currentTokens = math.Float64frombits(currentBits)

		if currentTokens < 1.0 {
			waitTime := time.Duration((1.0 - currentTokens) / c.rateLimit * float64(time.Second))
			if c.logger != nil {
				c.logger.Debug("api_client", "rate_limit_wait", fmt.Sprintf("no_tokens_available, wait_ms=%d", waitTime.Milliseconds()))
			}
			time.Sleep(waitTime)
		}
	}
}

// refillTokens refills the token bucket based on elapsed time
func (c *Client) refillTokens() {
	c.refillMu.Lock()
	defer c.refillMu.Unlock()

	now := time.Now().UnixNano()
	lastRefill := atomic.LoadInt64(&c.lastRefill)
	elapsed := float64(now-lastRefill) / 1e9

	if elapsed < 0.001 {
		return
	}

	tokensToAdd := elapsed * c.rateLimit
	currentBits := atomic.LoadUint64(&c.tokens)
	currentTokens := math.Float64frombits(currentBits)
	newTokens := math.Min(currentTokens+tokensToAdd, c.burstLimit)

	atomic.StoreUint64(&c.tokens, math.Float64bits(newTokens))
	atomic.StoreInt64(&c.lastRefill, now)

	if c.logger != nil {
		c.logger.Debug("api_client", "rate_limit_refill", fmt.Sprintf("added=%.2f, total=%.2f", tokensToAdd, newTokens))
	}
}
