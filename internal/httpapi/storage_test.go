package httpapi_test

import (
	"sync"
	"testing"
	"time"

	"github.com/joukojo/go-mail-testserver/internal/httpapi"
)

func TestNewStorage(t *testing.T) {
	s := httpapi.NewStorage()
	if s == nil {
		t.Fatal("NewStorage() returned nil")
	}

	messages := s.List()
	if len(messages) != 0 {
		t.Errorf("expected empty storage, got %d messages", len(messages))
	}
}

func TestStorage_Add(t *testing.T) {
	s := httpapi.NewStorage()

	msg := &httpapi.Message{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
		Raw:     []byte("raw email data"),
	}

	id := s.Add(msg)
	if id != 1 {
		t.Errorf("expected first ID to be 1, got %d", id)
	}

	if msg.ID != 1 {
		t.Errorf("expected message ID to be set to 1, got %d", msg.ID)
	}

	if msg.CreatedAt == "" {
		t.Error("expected CreatedAt to be set")
	}

	// Verify timestamp format
	_, err := time.Parse(time.RFC3339, msg.CreatedAt)
	if err != nil {
		t.Errorf("CreatedAt has invalid format: %v", err)
	}
}

func TestStorage_AddMultiple(t *testing.T) {
	s := httpapi.NewStorage()

	for i := 1; i <= 5; i++ {
		msg := &httpapi.Message{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test",
			Body:    "Test",
		}

		id := s.Add(msg)
		if id != i {
			t.Errorf("expected ID %d, got %d", i, id)
		}
	}

	messages := s.List()
	if len(messages) != 5 {
		t.Errorf("expected 5 messages, got %d", len(messages))
	}
}

func TestStorage_Get(t *testing.T) {
	s := httpapi.NewStorage()

	original := &httpapi.Message{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com", "another@example.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
		Raw:     []byte("raw email data"),
	}

	id := s.Add(original)

	retrieved, exists := s.Get(id)
	if !exists {
		t.Fatal("message not found")
	}

	if retrieved.ID != id {
		t.Errorf("expected ID %d, got %d", id, retrieved.ID)
	}
	if retrieved.From != original.From {
		t.Errorf("expected From %q, got %q", original.From, retrieved.From)
	}
	if len(retrieved.To) != len(original.To) {
		t.Errorf("expected %d recipients, got %d", len(original.To), len(retrieved.To))
	}
	if retrieved.Subject != original.Subject {
		t.Errorf("expected Subject %q, got %q", original.Subject, retrieved.Subject)
	}
	if retrieved.Body != original.Body {
		t.Errorf("expected Body %q, got %q", original.Body, retrieved.Body)
	}
	if string(retrieved.Raw) != string(original.Raw) {
		t.Errorf("expected Raw %q, got %q", original.Raw, retrieved.Raw)
	}
}

func TestStorage_GetNonExistent(t *testing.T) {
	s := httpapi.NewStorage()

	_, exists := s.Get(999)
	if exists {
		t.Error("expected non-existent message to return false")
	}
}

func TestStorage_GetReturnsImmutableCopy(t *testing.T) {
	s := httpapi.NewStorage()

	msg := &httpapi.Message{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Original Subject",
		Body:    "Original Body",
		Raw:     []byte("original raw data"),
	}

	id := s.Add(msg)

	retrieved, _ := s.Get(id)
	retrieved.Subject = "Modified Subject"
	retrieved.To[0] = "modified@example.com"
	retrieved.Raw[0] = 'X'

	// Get again and verify original values
	original, _ := s.Get(id)
	if original.Subject != "Original Subject" {
		t.Errorf("expected original Subject, got %q", original.Subject)
	}
	if original.To[0] != "recipient@example.com" {
		t.Errorf("expected original To, got %q", original.To[0])
	}
	if original.Raw[0] != 'o' {
		t.Errorf("expected original Raw, got %c", original.Raw[0])
	}
}

func TestStorage_List(t *testing.T) {
	s := httpapi.NewStorage()

	messages := []httpapi.Message{
		{From: "sender1@example.com", To: []string{"recipient1@example.com"}, Subject: "Subject 1", Body: "Body 1"},
		{From: "sender2@example.com", To: []string{"recipient2@example.com"}, Subject: "Subject 2", Body: "Body 2"},
		{From: "sender3@example.com", To: []string{"recipient3@example.com"}, Subject: "Subject 3", Body: "Body 3"},
	}

	for i := range messages {
		s.Add(&messages[i])
	}

	list := s.List()
	if len(list) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(list))
	}

	// Verify all messages are present
	for _, msg := range messages {
		found := false
		for _, listedMsg := range list {
			if listedMsg.From == msg.From && listedMsg.Subject == msg.Subject {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("message from %q not found in list", msg.From)
		}
	}
}

func TestStorage_ListReturnsImmutableCopy(t *testing.T) {
	s := httpapi.NewStorage()

	msg := &httpapi.Message{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Original Subject",
		Body:    "Original Body",
	}

	s.Add(msg)

	list := s.List()
	list[0].Subject = "Modified Subject"
	list[0].To[0] = "modified@example.com"

	// Get list again and verify original values
	list2 := s.List()
	if list2[0].Subject != "Original Subject" {
		t.Errorf("expected original Subject, got %q", list2[0].Subject)
	}
	if list2[0].To[0] != "recipient@example.com" {
		t.Errorf("expected original To, got %q", list2[0].To[0])
	}
}

func TestStorage_ListDoesNotExposeRawData(t *testing.T) {
	s := httpapi.NewStorage()

	msg := &httpapi.Message{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test",
		Body:    "Test",
		Raw:     []byte("secret raw data"),
	}

	s.Add(msg)

	list := s.List()
	if len(list[0].Raw) != 0 {
		t.Error("List() should not expose Raw data")
	}
}

func TestStorage_Clear(t *testing.T) {
	s := httpapi.NewStorage()

	for i := 0; i < 5; i++ {
		s.Add(&httpapi.Message{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test",
			Body:    "Test",
		})
	}

	s.Clear()

	messages := s.List()
	if len(messages) != 0 {
		t.Errorf("expected empty storage after Clear(), got %d messages", len(messages))
	}

	// Verify ID counter is reset
	msg := &httpapi.Message{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test",
		Body:    "Test",
	}
	id := s.Add(msg)
	if id != 1 {
		t.Errorf("expected ID to reset to 1 after Clear(), got %d", id)
	}
}

// Concurrency Tests

func TestStorage_ConcurrentAdd(t *testing.T) {
	s := httpapi.NewStorage()
	numGoroutines := 100
	messagesPerGoroutine := 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				msg := &httpapi.Message{
					From:    "sender@example.com",
					To:      []string{"recipient@example.com"},
					Subject: "Test",
					Body:    "Test",
				}
				s.Add(msg)
			}
		}(i)
	}

	wg.Wait()

	messages := s.List()
	expectedCount := numGoroutines * messagesPerGoroutine
	if len(messages) != expectedCount {
		t.Errorf("expected %d messages, got %d", expectedCount, len(messages))
	}

	// Verify all IDs are unique
	idMap := make(map[int]bool)
	for _, msg := range messages {
		if idMap[msg.ID] {
			t.Errorf("duplicate ID found: %d", msg.ID)
		}
		idMap[msg.ID] = true
	}
}

func TestStorage_ConcurrentRead(t *testing.T) {
	s := httpapi.NewStorage()

	// Add some messages first
	for i := 0; i < 10; i++ {
		s.Add(&httpapi.Message{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test",
			Body:    "Test",
		})
	}

	numGoroutines := 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 1; j <= 10; j++ {
				msg, exists := s.Get(j)
				if !exists {
					errors <- nil
					return
				}
				if msg.ID != j {
					errors <- nil
					return
				}
			}
		}()
	}

	wg.Wait()
	close(errors)

	for range errors {
		t.Error("concurrent read failed")
	}
}

func TestStorage_ConcurrentList(t *testing.T) {
	s := httpapi.NewStorage()

	// Add some messages
	for i := 0; i < 20; i++ {
		s.Add(&httpapi.Message{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test",
			Body:    "Test",
		})
	}

	numGoroutines := 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			messages := s.List()
			if len(messages) != 20 {
				t.Errorf("expected 20 messages, got %d", len(messages))
			}
		}()
	}

	wg.Wait()
}

func TestStorage_ConcurrentAddAndRead(t *testing.T) {
	s := httpapi.NewStorage()
	numWriters := 50
	numReaders := 50
	duration := 100 * time.Millisecond

	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Writers
	wg.Add(numWriters)
	for i := 0; i < numWriters; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					s.Add(&httpapi.Message{
						From:    "sender@example.com",
						To:      []string{"recipient@example.com"},
						Subject: "Test",
						Body:    "Test",
					})
				}
			}
		}()
	}

	// Readers
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					s.List()
				}
			}
		}()
	}

	time.Sleep(duration)
	close(stop)
	wg.Wait()

	// Verify data integrity
	messages := s.List()
	idMap := make(map[int]bool)
	for _, msg := range messages {
		if idMap[msg.ID] {
			t.Errorf("duplicate ID found: %d", msg.ID)
		}
		idMap[msg.ID] = true

		if msg.From != "sender@example.com" {
			t.Errorf("corrupted From field: %q", msg.From)
		}
	}
}

func TestStorage_ConcurrentAddAndGet(t *testing.T) {
	s := httpapi.NewStorage()
	numGoroutines := 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Mix of Add and Get operations
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()

			if idx%2 == 0 {
				// Add messages
				s.Add(&httpapi.Message{
					From:    "sender@example.com",
					To:      []string{"recipient@example.com"},
					Subject: "Test",
					Body:    "Test",
				})
			} else {
				// Try to get messages
				s.Get(idx / 2)
			}
		}(i)
	}

	wg.Wait()

	// Verify no corruption
	messages := s.List()
	if len(messages) == 0 {
		t.Error("expected messages to be added")
	}
}

func TestStorage_ConcurrentClear(t *testing.T) {
	s := httpapi.NewStorage()

	// Add initial messages
	for i := 0; i < 10; i++ {
		s.Add(&httpapi.Message{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test",
			Body:    "Test",
		})
	}

	var wg sync.WaitGroup
	numGoroutines := 10
	wg.Add(numGoroutines)

	// Multiple goroutines calling Clear simultaneously
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			s.Clear()
		}()
	}

	wg.Wait()

	// Verify storage is empty
	messages := s.List()
	if len(messages) != 0 {
		t.Errorf("expected empty storage, got %d messages", len(messages))
	}

	// Verify next ID is 1
	msg := &httpapi.Message{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test",
		Body:    "Test",
	}
	id := s.Add(msg)
	if id != 1 {
		t.Errorf("expected ID to be 1, got %d", id)
	}
}

func TestStorage_ConcurrentMixedOperations(t *testing.T) {
	s := httpapi.NewStorage()
	duration := 200 * time.Millisecond
	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Adders
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				s.Add(&httpapi.Message{
					From:    "sender@example.com",
					To:      []string{"recipient@example.com"},
					Subject: "Test",
					Body:    "Test",
				})
			}
		}
	}()

	// Listers
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				s.List()
			}
		}
	}()

	// Getters
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				s.Get(1)
			}
		}
	}()

	// Periodic clearers
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				s.Clear()
			}
		}
	}()

	time.Sleep(duration)
	close(stop)
	wg.Wait()

	// Just verify we didn't crash
	s.List()
}

func BenchmarkStorage_Add(b *testing.B) {
	s := httpapi.NewStorage()
	msg := &httpapi.Message{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test",
		Body:    "Test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(msg)
	}
}

func BenchmarkStorage_Get(b *testing.B) {
	s := httpapi.NewStorage()
	msg := &httpapi.Message{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test",
		Body:    "Test",
	}
	id := s.Add(msg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Get(id)
	}
}

func BenchmarkStorage_List(b *testing.B) {
	s := httpapi.NewStorage()
	for i := 0; i < 100; i++ {
		s.Add(&httpapi.Message{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test",
			Body:    "Test",
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.List()
	}
}

func BenchmarkStorage_ConcurrentAdd(b *testing.B) {
	s := httpapi.NewStorage()
	msg := &httpapi.Message{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test",
		Body:    "Test",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Add(msg)
		}
	})
}

func BenchmarkStorage_ConcurrentRead(b *testing.B) {
	s := httpapi.NewStorage()
	for i := 0; i < 100; i++ {
		s.Add(&httpapi.Message{
			From:    "sender@example.com",
			To:      []string{"recipient@example.com"},
			Subject: "Test",
			Body:    "Test",
		})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.List()
		}
	})
}
