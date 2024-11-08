package commands

import (
	"sync"

	"github.com/google/generative-ai-go/genai"
)

// ChatSession represents a single chat session with history
type ChatSession struct {
	History []*genai.Content
}

// SessionManager manages chat sessions for different users
type SessionManager struct {
	mu       sync.RWMutex
	sessions map[int64]*ChatSession
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[int64]*ChatSession),
	}
}

// GetOrCreate gets an existing session or creates a new one
func (sm *SessionManager) GetOrCreate(chatID int64) *ChatSession {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if session, exists := sm.sessions[chatID]; exists {
		return session
	}

	session := &ChatSession{
		History: make([]*genai.Content, 0),
	}
	sm.sessions[chatID] = session
	return session
}

// Clear clears the chat history for a given chat ID
func (sm *SessionManager) Clear(chatID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if session, exists := sm.sessions[chatID]; exists {
		session.History = make([]*genai.Content, 0)
	}
}
