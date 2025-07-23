package local

import (
	"fmt"
	"sync"
	"testing"
)

func TestFSMBasicOperations(t *testing.T) {
	fsm := NewFSM()

	t.Run("CreateUser", func(t *testing.T) {
		success := fsm.CreateUser(123, "initial")
		if !success {
			t.Error("Expected CreateUser to succeed")
		}

		success = fsm.CreateUser(123, "duplicate")
		if success {
			t.Error("Expected CreateUser to fail for duplicate user")
		}
	})

	t.Run("GetStateUser", func(t *testing.T) {
		state, exists := fsm.GetStateUser(123)
		if !exists {
			t.Error("Expected user to exist")
		}
		if state != "initial" {
			t.Errorf("Expected state 'initial', got '%s'", state)
		}

		_, exists = fsm.GetStateUser(999)
		if exists {
			t.Error("Expected non-existent user to return false")
		}
	})

	t.Run("SetStateUser", func(t *testing.T) {
		success := fsm.SetStateUser(123, "updated")
		if !success {
			t.Error("Expected SetStateUser to succeed")
		}

		state, _ := fsm.GetStateUser(123)
		if state != "updated" {
			t.Errorf("Expected state 'updated', got '%s'", state)
		}

		success = fsm.SetStateUser(999, "nonexistent")
		if success {
			t.Error("Expected SetStateUser to fail for non-existent user")
		}
	})

	t.Run("SetAndGetValue", func(t *testing.T) {
		success := fsm.SetValue(123, "key1", "value1")
		if !success {
			t.Error("Expected SetValue to succeed")
		}

		value, exists := fsm.GetValue(123, "key1")
		if !exists {
			t.Error("Expected value to exist")
		}
		if value != "value1" {
			t.Errorf("Expected value 'value1', got '%s'", value)
		}

		_, exists = fsm.GetValue(123, "nonexistent")
		if exists {
			t.Error("Expected non-existent key to return false")
		}

		success = fsm.SetValue(999, "key", "value")
		if success {
			t.Error("Expected SetValue to fail for non-existent user")
		}

		_, exists = fsm.GetValue(999, "key")
		if exists {
			t.Error("Expected GetValue to fail for non-existent user")
		}
	})
}

func TestFSMConcurrentAccess(t *testing.T) {
	fsm := NewFSM()
	userID := int64(100)

	fsm.CreateUser(userID, "initial")

	t.Run("ConcurrentStateUpdates", func(t *testing.T) {
		const numGoroutines = 50
		const numOperations = 100

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(workerID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					state := fmt.Sprintf("state-%d-%d", workerID, j)
					fsm.SetStateUser(userID, state)
					fsm.GetStateUser(userID)
				}
			}(i)
		}

		wg.Wait()

		state, exists := fsm.GetStateUser(userID)
		if !exists {
			t.Error("Expected user to still exist after concurrent updates")
		}
		if state == "" {
			t.Error("Expected state to not be empty")
		}
	})

	t.Run("ConcurrentValueOperations", func(t *testing.T) {
		const numGoroutines = 30
		const numOperations = 50

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(workerID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					key := fmt.Sprintf("key-%d", workerID)
					value := fmt.Sprintf("value-%d-%d", workerID, j)
					
					fsm.SetValue(userID, key, value)
					fsm.GetValue(userID, key)
				}
			}(i)
		}

		wg.Wait()

		fsm.SetValue(userID, "test", "final")
		value, exists := fsm.GetValue(userID, "test")
		if !exists || value != "final" {
			t.Error("Expected to be able to set and get value after concurrent operations")
		}
	})

	t.Run("ConcurrentUserCreation", func(t *testing.T) {
		const numGoroutines = 20
		var wg sync.WaitGroup
		var successCount int32
		var mu sync.Mutex

		wg.Add(numGoroutines)

		targetUserID := int64(200)
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				if fsm.CreateUser(targetUserID, "concurrent") {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}()
		}

		wg.Wait()

		if successCount != 1 {
			t.Errorf("Expected exactly 1 successful creation, got %d", successCount)
		}

		_, exists := fsm.GetStateUser(targetUserID)
		if !exists {
			t.Error("Expected user to exist after concurrent creation attempts")
		}
	})
}

func TestFSMUserLifecycle(t *testing.T) {
	fsm := NewFSM()
	userID := int64(300)

	t.Run("CompleteLifecycle", func(t *testing.T) {
		_, exists := fsm.GetStateUser(userID)
		if exists {
			t.Error("Expected user to not exist initially")
		}

		success := fsm.CreateUser(userID, "registered")
		if !success {
			t.Fatal("Expected user creation to succeed")
		}

		state, exists := fsm.GetStateUser(userID)
		if !exists || state != "registered" {
			t.Error("Expected user to exist with 'registered' state")
		}

		testData := map[string]string{
			"name":     "John Doe",
			"email":    "john@example.com",
			"settings": "dark_mode",
		}

		for key, value := range testData {
			success := fsm.SetValue(userID, key, value)
			if !success {
				t.Errorf("Expected SetValue to succeed for key '%s'", key)
			}
		}

		for key, expectedValue := range testData {
			value, exists := fsm.GetValue(userID, key)
			if !exists {
				t.Errorf("Expected value to exist for key '%s'", key)
			}
			if value != expectedValue {
				t.Errorf("Expected value '%s' for key '%s', got '%s'", expectedValue, key, value)
			}
		}

		states := []string{"verified", "active", "premium", "inactive"}
		for _, newState := range states {
			success := fsm.SetStateUser(userID, newState)
			if !success {
				t.Errorf("Expected state update to '%s' to succeed", newState)
			}

			currentState, _ := fsm.GetStateUser(userID)
			if currentState != newState {
				t.Errorf("Expected current state to be '%s', got '%s'", newState, currentState)
			}
		}

		for key, expectedValue := range testData {
			value, exists := fsm.GetValue(userID, key)
			if !exists || value != expectedValue {
				t.Errorf("Expected value '%s' for key '%s' to persist after state changes", expectedValue, key)
			}
		}
	})
}

func TestFSMEdgeCases(t *testing.T) {
	fsm := NewFSM()

	t.Run("EmptyState", func(t *testing.T) {
		userID := int64(400)
		success := fsm.CreateUser(userID, "")
		if !success {
			t.Error("Expected CreateUser to succeed with empty state")
		}

		state, exists := fsm.GetStateUser(userID)
		if !exists || state != "" {
			t.Error("Expected user to exist with empty state")
		}
	})

	t.Run("EmptyValues", func(t *testing.T) {
		userID := int64(401)
		fsm.CreateUser(userID, "test")

		success := fsm.SetValue(userID, "empty", "")
		if !success {
			t.Error("Expected SetValue to succeed with empty value")
		}

		value, exists := fsm.GetValue(userID, "empty")
		if !exists || value != "" {
			t.Error("Expected empty value to be stored and retrieved")
		}

		success = fsm.SetValue(userID, "", "value")
		if !success {
			t.Error("Expected SetValue to succeed with empty key")
		}

		value, exists = fsm.GetValue(userID, "")
		if !exists || value != "value" {
			t.Error("Expected value with empty key to be stored and retrieved")
		}
	})

	t.Run("LargeUserIDs", func(t *testing.T) {
		largeID := int64(9223372036854775807)
		success := fsm.CreateUser(largeID, "large")
		if !success {
			t.Error("Expected CreateUser to succeed with large user ID")
		}

		state, exists := fsm.GetStateUser(largeID)
		if !exists || state != "large" {
			t.Error("Expected user with large ID to exist")
		}
	})

	t.Run("NegativeUserIDs", func(t *testing.T) {
		negativeID := int64(-123)
		success := fsm.CreateUser(negativeID, "negative")
		if !success {
			t.Error("Expected CreateUser to succeed with negative user ID")
		}

		state, exists := fsm.GetStateUser(negativeID)
		if !exists || state != "negative" {
			t.Error("Expected user with negative ID to exist")
		}
	})
}

// Benchmark tests
func BenchmarkFSMOperations(b *testing.B) {
	fsm := NewFSM()
	userID := int64(500)
	fsm.CreateUser(userID, "benchmark")

	b.Run("GetStateUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fsm.GetStateUser(userID)
		}
	})

	b.Run("SetStateUser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			state := fmt.Sprintf("state-%d", i)
			fsm.SetStateUser(userID, state)
		}
	})

	b.Run("SetValue", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key-%d", i%100)
			value := fmt.Sprintf("value-%d", i)
			fsm.SetValue(userID, key, value)
		}
	})

	b.Run("GetValue", func(b *testing.B) {
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("key-%d", i)
			fsm.SetValue(userID, key, "value")
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key-%d", i%100)
			fsm.GetValue(userID, key)
		}
	})
}

func BenchmarkFSMConcurrentLoad(b *testing.B) {
	fsm := NewFSM()
	
	const numUsers = 1000
	for i := 0; i < numUsers; i++ {
		fsm.CreateUser(int64(i), "benchmark")
	}

	b.Run("ConcurrentMixedOperations", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			userID := int64(0)
			for pb.Next() {
				userID = (userID + 1) % numUsers
				
					fsm.GetStateUser(userID)
				fsm.SetValue(userID, "bench", "value")
				fsm.GetValue(userID, "bench")
			}
		})
	})
}