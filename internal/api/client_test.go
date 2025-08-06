package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexMayka/go-max-sdk/internal"
	"github.com/AlexMayka/go-max-sdk/internal/core"
	"github.com/AlexMayka/go-max-sdk/types/models"
	reqMsg "github.com/AlexMayka/go-max-sdk/types/requests/messages"
	resMsg "github.com/AlexMayka/go-max-sdk/types/responses/messages"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestClientBasicOperations(t *testing.T) {
	t.Run("SuccessfulRequest", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Expected POST method, got %s", r.Method)
			}

			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected application/json content type, got %s", r.Header.Get("Content-Type"))
			}

			if !strings.Contains(r.URL.RawQuery, "access_token=test_token") {
				t.Error("Expected access_token in query parameters")
			}

			response := resMsg.Send{
				Message: models.Message{
					Timestamp: 1234567890,
					Body: &models.MessageBody{
						Text: "Hello, World!",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		originalHost := Host
		originalScheme := Scheme
		defer func() {
			Host = originalHost
			Scheme = originalScheme
		}()

		serverURL := strings.TrimPrefix(server.URL, "http://")
		Host = serverURL
		Scheme = "http"

		client := NewClient("test_token", internal.NewNoopLogger(), 5*time.Second, 10, 5, 3, 100*time.Millisecond)

		request := reqMsg.Send{
			NewMessageBody: models.NewMessageBody{
				Text: stringPtr("Hello, World!"),
			},
			UserID: int64Ptr(12345),
		}

		ctx := context.Background()
		response, err := client.Call(ctx, core.SendMsg, request)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		sendResponse, ok := response.(*resMsg.Send)
		if !ok {
			t.Fatalf("Expected *resMsg.Send, got %T", response)
		}

		if sendResponse.Message.Body == nil || sendResponse.Message.Body.Text != "Hello, World!" {
			t.Errorf("Expected message text 'Hello, World!', got %v", sendResponse.Message.Body)
		}
	})

	t.Run("HTTPError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
		}))
		defer server.Close()

		originalHost := Host
		originalScheme := Scheme
		defer func() {
			Host = originalHost
			Scheme = originalScheme
		}()

		serverURL := strings.TrimPrefix(server.URL, "http://")
		Host = serverURL
		Scheme = "http"

		client := NewClient("test_token", internal.NewNoopLogger(), 5*time.Second, 10, 5, 0, 100*time.Millisecond)

		request := reqMsg.Send{
			NewMessageBody: models.NewMessageBody{
				Text: stringPtr("Hello"),
			},
		}

		ctx := context.Background()
		_, err := client.Call(ctx, core.SendMsg, request)

		if err == nil {
			t.Fatal("Expected error for HTTP 400, got nil")
		}
	})
}

func TestClientRateLimiting(t *testing.T) {
	t.Run("RateLimitingWorks", func(t *testing.T) {
		requestCount := 0
		var mu sync.Mutex

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			requestCount++
			mu.Unlock()

			response := resMsg.Send{
				Message: models.Message{
					Timestamp: 1234567890,
					Body: &models.MessageBody{
						Text: fmt.Sprintf("msg_%d", requestCount),
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		originalHost := Host
		originalScheme := Scheme
		defer func() {
			Host = originalHost
			Scheme = originalScheme
		}()

		serverURL := strings.TrimPrefix(server.URL, "http://")
		Host = serverURL
		Scheme = "http"

		// Настраиваем клиент с rate limiting: 2 запроса в секунду
		client := NewClient("test_token", internal.NewNoopLogger(), 5*time.Second, 2.0, 2.0, 0, 100*time.Millisecond)

		request := reqMsg.Send{
			NewMessageBody: models.NewMessageBody{
				Text: stringPtr("Test"),
			},
		}

		ctx := context.Background()
		start := time.Now()

		for i := 0; i < 4; i++ {
			_, err := client.Call(ctx, core.SendMsg, request)
			if err != nil {
				t.Fatalf("Request %d failed: %v", i, err)
			}
		}

		elapsed := time.Since(start)

		// С rate limit 2 req/sec, 4 запроса должны занять минимум 1 секунду
		if elapsed < 800*time.Millisecond {
			t.Errorf("Expected rate limiting to slow down requests, but took only %v", elapsed)
		}

		if requestCount != 4 {
			t.Errorf("Expected 4 requests, got %d", requestCount)
		}
	})
}

func TestClientRetryLogic(t *testing.T) {
	t.Run("RetriesOnFailure", func(t *testing.T) {
		attemptCount := 0
		var mu sync.Mutex

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			attemptCount++
			currentAttempt := attemptCount
			mu.Unlock()

			// Первые 2 запроса возвращают ошибку, третий - успех
			if currentAttempt <= 2 {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Server Error"))
				return
			}

			response := resMsg.Send{
				Message: models.Message{
					Timestamp: 1234567890,
					Body: &models.MessageBody{
						Text: "success_message",
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		originalHost := Host
		originalScheme := Scheme
		defer func() {
			Host = originalHost
			Scheme = originalScheme
		}()

		serverURL := strings.TrimPrefix(server.URL, "http://")
		Host = serverURL
		Scheme = "http"

		client := NewClient("test_token", internal.NewNoopLogger(), 5*time.Second, 10, 0, 3, 10*time.Millisecond)

		request := reqMsg.Send{
			NewMessageBody: models.NewMessageBody{
				Text: stringPtr("Retry test"),
			},
		}

		ctx := context.Background()
		response, err := client.Call(ctx, core.SendMsg, request)

		if err != nil {
			t.Fatalf("Expected success after retries, got error: %v", err)
		}

		sendResponse, ok := response.(*resMsg.Send)
		if !ok || sendResponse.Message.Body.Text != "success_message" {
			t.Errorf("Expected successful response, got %v", response)
		}

		if attemptCount != 3 {
			t.Errorf("Expected 3 attempts, got %d", attemptCount)
		}
	})
}

func TestRequestParsing(t *testing.T) {
	t.Run("ParsePathParams", func(t *testing.T) {
		client := &Client{token: "test"}

		cfg := EndpointConfigs[GetMsgByID]

		request := reqMsg.GetByID{
			MessageID: "test_message_123",
		}

		pathParams, queryParams, jsonBody, err := client.parseRequest(request, cfg)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if pathParams["messageId"] != "test_message_123" {
			t.Errorf("Expected messageId 'test_message_123', got '%s'", pathParams["messageId"])
		}

		if len(queryParams) != 0 {
			t.Errorf("Expected no query params, got %v", queryParams)
		}

		if jsonBody == nil {
			t.Error("Expected jsonBody to be non-nil")
		}
	})

	t.Run("InvalidRequestType", func(t *testing.T) {
		client := &Client{token: "test"}
		cfg := EndpointConfigs[core.SendMsg]

		wrongRequest := reqMsg.GetByID{MessageID: "test"}

		_, _, _, err := client.parseRequest(wrongRequest, cfg)

		if err == nil {
			t.Fatal("Expected error for invalid request type, got nil")
		}

		if !strings.Contains(err.Error(), "invalid request type") {
			t.Errorf("Expected 'invalid request type' error, got: %v", err)
		}
	})
}

func TestURLBuilding(t *testing.T) {
	t.Run("BuildURLWithPathAndQuery", func(t *testing.T) {
		client := &Client{token: "test_token_123"}

		pathParams := map[string]string{
			"chatId":    "12345",
			"messageId": "msg_67890",
		}

		queryParams := map[string]string{
			"limit":  "50",
			"offset": "100",
		}

		url := client.buildURL("/chats/{chatId}/messages/{messageId}", pathParams, queryParams)

		expectedParts := []string{
			"botapi.max.ru/chats/12345/messages/msg_67890",
			"access_token=test_token_123",
			"limit=50",
			"offset=100",
		}

		for _, part := range expectedParts {
			if !strings.Contains(url, part) {
				t.Errorf("Expected URL to contain '%s', got: %s", part, url)
			}
		}
	})

	t.Run("BuildURLWithoutToken", func(t *testing.T) {
		client := &Client{token: ""}

		url := client.buildURL("/test", map[string]string{}, map[string]string{})

		if strings.Contains(url, "access_token") {
			t.Errorf("Expected no access_token in URL without token, got: %s", url)
		}
	})
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}
