package httpapi

import (
	"sync"
	"time"
)

// Message represents an email message
type Message struct {
	ID        int      `json:"id"`
	From      string   `json:"from"`
	To        []string `json:"to"`
	Subject   string   `json:"subject"`
	Body      string   `json:"body"`
	CreatedAt string   `json:"createdAt"`
	Raw       []byte   `json:"-"` // RFC822 raw bytes, not exposed in JSON
}

// Storage manages email messages with thread-safe operations
type Storage struct {
	mu       sync.RWMutex
	messages map[int]*Message
	nextID   int
}

func NewStorage() *Storage {
	return &Storage{
		messages: make(map[int]*Message),
		nextID:   1,
	}
}

// Add stores a new message and returns its ID
func (s *Storage) Add(msg *Message) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg.ID = s.nextID
	msg.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	s.messages[msg.ID] = msg
	s.nextID++
	return msg.ID
}

// List returns a copy of all messages (without raw data)
func (s *Storage) List() []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Message, 0, len(s.messages))
	for _, msg := range s.messages {
		// Copy message without raw bytes
		result = append(result, Message{
			ID:        msg.ID,
			From:      msg.From,
			To:        append([]string(nil), msg.To...),
			Subject:   msg.Subject,
			Body:      msg.Body,
			CreatedAt: msg.CreatedAt,
		})
	}
	return result
}

// Get retrieves a message by ID
func (s *Storage) Get(id int) (*Message, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msg, exists := s.messages[id]
	if !exists {
		return nil, false
	}

	// Return a copy
	msgCopy := &Message{
		ID:        msg.ID,
		From:      msg.From,
		To:        append([]string(nil), msg.To...),
		Subject:   msg.Subject,
		Body:      msg.Body,
		CreatedAt: msg.CreatedAt,
		Raw:       append([]byte(nil), msg.Raw...),
	}
	return msgCopy, true
}

// Clear removes all messages
func (s *Storage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages = make(map[int]*Message)
	s.nextID = 1
}
