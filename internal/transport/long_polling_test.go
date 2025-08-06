package transport

import (
	"context"
	"errors"
	"github.com/AlexMayka/go-max-sdk/internal"
	"github.com/AlexMayka/go-max-sdk/internal/core"
	"sync"
	"testing"
	"time"

	"github.com/AlexMayka/go-max-sdk/types/models"
	req "github.com/AlexMayka/go-max-sdk/types/requests/subscriptions"
	res "github.com/AlexMayka/go-max-sdk/types/responses/subscriptions"
)

type mockAPIClient struct {
	responses []interface{}
	errors    []error
	callCount int
	mu        sync.Mutex

	// Для проверки параметров
	lastEndpoint core.Endpoint
	lastRequest  interface{}
}

func (m *mockAPIClient) Call(ctx context.Context, endpoint core.Endpoint, request interface{}) (interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastEndpoint = endpoint
	m.lastRequest = request

	if m.callCount < len(m.errors) && m.errors[m.callCount] != nil {
		err := m.errors[m.callCount]
		m.callCount++
		return nil, err
	}

	if m.callCount < len(m.responses) {
		response := m.responses[m.callCount]
		m.callCount++
		return response, nil
	}

	// Возвращаем пустой ответ если нет больше мок данных
	return res.GetUpdates{Updates: []models.Update{}}, nil
}

func TestLongPollingBasic(t *testing.T) {
	t.Run("SuccessfulPolling", func(t *testing.T) {
		updates := []models.Update{
			{UpdateType: models.UpdateTypeMessageCreated},
			{UpdateType: models.UpdateTypeMessageCallback},
		}

		mockClient := &mockAPIClient{
			responses: []interface{}{
				res.GetUpdates{Updates: updates},
			},
		}

		transport := NewLongPolling(30, 100, 10, time.Second, 3, mockClient, internal.NewNoopLogger())

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		ch, err := transport.Start(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Проверяем что получили updates
		receivedCount := 0
		timeout := time.After(200 * time.Millisecond)

	loop:
		for {
			select {
			case update, ok := <-ch:
				if !ok {
					break loop
				}
				if update != nil {
					receivedCount++
				}
			case <-timeout:
				break loop
			}
		}

		if receivedCount != len(updates) {
			t.Errorf("Expected %d updates, got %d", len(updates), receivedCount)
		}

		transport.Stop()
	})
}

func TestLongPollingRetry(t *testing.T) {
	t.Run("RetryOnError", func(t *testing.T) {
		mockClient := &mockAPIClient{
			errors: []error{
				errors.New("network error"),
				errors.New("timeout error"),
				nil, // третий запрос успешный
			},
			responses: []interface{}{
				nil, nil, // первые два ответа игнорируются из-за ошибок
				res.GetUpdates{Updates: []models.Update{}},
			},
		}

		transport := NewLongPolling(30, 100, 10, 10*time.Millisecond, 5, mockClient, internal.NewNoopLogger())

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		ch, err := transport.Start(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Ждем немного чтобы произошли retry
		time.Sleep(50 * time.Millisecond)

		// Проверяем что было несколько вызовов
		if mockClient.callCount < 3 {
			t.Errorf("Expected at least 3 calls (with retries), got %d", mockClient.callCount)
		}

		transport.Stop()

		// Проверяем что канал закрылся
		select {
		case _, ok := <-ch:
			if ok {
				t.Error("Expected channel to be closed")
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Channel should be closed quickly after Stop()")
		}
	})
}

func TestLongPollingGracefulShutdown(t *testing.T) {
	t.Run("StopClosesChannel", func(t *testing.T) {
		mockClient := &mockAPIClient{
			responses: []interface{}{
				res.GetUpdates{Updates: []models.Update{}},
			},
		}

		transport := NewLongPolling(30, 100, 10, time.Second, 3, mockClient, internal.NewNoopLogger())

		ctx := context.Background()
		ch, err := transport.Start(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Останавливаем transport
		err = transport.Stop()
		if err != nil {
			t.Errorf("Expected no error on stop, got %v", err)
		}

		// Проверяем что канал закрылся
		select {
		case _, ok := <-ch:
			if ok {
				t.Error("Expected channel to be closed after Stop()")
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Channel should be closed quickly after Stop()")
		}
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		mockClient := &mockAPIClient{
			responses: []interface{}{
				res.GetUpdates{Updates: []models.Update{}},
			},
		}

		transport := NewLongPolling(30, 100, 10, time.Second, 3, mockClient, internal.NewNoopLogger())

		ctx, cancel := context.WithCancel(context.Background())
		ch, err := transport.Start(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Отменяем контекст
		cancel()

		// Проверяем что канал закрылся
		select {
		case _, ok := <-ch:
			if ok {
				t.Error("Expected channel to be closed after context cancellation")
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Channel should be closed quickly after context cancellation")
		}
	})
}

func TestLongPollingConfiguration(t *testing.T) {
	t.Run("CorrectRequestParameters", func(t *testing.T) {
		mockClient := &mockAPIClient{
			responses: []interface{}{
				res.GetUpdates{Updates: []models.Update{}},
			},
		}

		pollTimeout := 45
		pollLimit := 50

		transport := NewLongPolling(pollTimeout, pollLimit, 10, time.Second, 3, mockClient, internal.NewNoopLogger())

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		_, err := transport.Start(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Ждем первый запрос
		time.Sleep(20 * time.Millisecond)

		// Проверяем параметры запроса
		if mockClient.lastEndpoint != core.GetSubscribeUpdate {
			t.Errorf("Expected endpoint %v, got %v", core.GetSubscribeUpdate, mockClient.lastEndpoint)
		}

		if mockClient.lastRequest == nil {
			t.Fatal("Expected request to be set")
		}

		getUpdatesReq, ok := mockClient.lastRequest.(*req.GetUpdates)
		if !ok {
			t.Fatalf("Expected GetUpdates request, got %T", mockClient.lastRequest)
		}

		if getUpdatesReq.Timeout == nil || *getUpdatesReq.Timeout != pollTimeout {
			t.Errorf("Expected timeout %d, got %v", pollTimeout, getUpdatesReq.Timeout)
		}

		if getUpdatesReq.Limit == nil || *getUpdatesReq.Limit != pollLimit {
			t.Errorf("Expected limit %d, got %v", pollLimit, getUpdatesReq.Limit)
		}

		if getUpdatesReq.Types == nil || len(getUpdatesReq.Types) == 0 {
			t.Error("Expected Types to be set")
		}

		transport.Stop()
	})
}

func TestConvertUpdateType(t *testing.T) {
	t.Run("ConvertMultipleTypes", func(t *testing.T) {
		types := convertUpdateType(
			models.UpdateTypeMessageCreated,
			models.UpdateTypeMessageCallback,
			models.UpdateBotStarted,
		)

		expected := []string{
			string(models.UpdateTypeMessageCreated),
			string(models.UpdateTypeMessageCallback),
			string(models.UpdateBotStarted),
		}

		if len(types) != len(expected) {
			t.Errorf("Expected %d types, got %d", len(expected), len(types))
		}

		for i, expectedType := range expected {
			if i >= len(types) || types[i] != expectedType {
				t.Errorf("Expected type[%d] = %s, got %s", i, expectedType, types[i])
			}
		}
	})

	t.Run("ConvertEmptyTypes", func(t *testing.T) {
		types := convertUpdateType()

		if len(types) != 0 {
			t.Errorf("Expected empty slice, got %v", types)
		}
	})
}
