package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Email struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	Date    string   `json:"date"`
}

type Server struct {
	addr string

	storage  *Storage
	emailsMu sync.RWMutex
}

// New creates a new HTTP API server
func New(addr string, storage *Storage) *Server {
	return &Server{
		addr:    addr,
		storage: storage,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/api/v1/messages", s.handleEmails)
	mux.HandleFunc("/api/v1/messages/clear", s.handleClear)
	mux.HandleFunc("/api/v1/messages/{id}", s.handleEmail)
	mux.HandleFunc("/api/v1/messages/{id}/raw", s.handleRawEmail)
	// mux.HandleFunc("/health", s.handleHealth)

	return http.ListenAndServe(s.addr, mux)
}

// handleEmails returns all received emails
func (s *Server) handleEmails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.emailsMu.RLock()
	defer s.emailsMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.storage.List())
}

func (s *Server) handleClear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.emailsMu.Lock()
	defer s.emailsMu.Unlock()

	s.storage.Clear()

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleEmail(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Provision ID: ", r.PathValue("id"))

	s.emailsMu.RLock()
	defer s.emailsMu.RUnlock()

	idStr := r.PathValue("id")
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	msg, exists := s.storage.Get(id)
	if !exists {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) handleRawEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.emailsMu.RLock()
	defer s.emailsMu.RUnlock()

	idStr := r.PathValue("id")
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	msg, exists := s.storage.Get(id)
	if !exists {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(msg.Raw)
}
