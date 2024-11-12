package sessions

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Username string
	IsAdmin  bool
	Expiry   time.Time
}

type SessionsStore struct {
	store map[string]Session
	mu    sync.Mutex
}

var Sessions SessionsStore = SessionsStore{store: map[string]Session{}}

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func (s *SessionsStore) Add(sessionToken string, session Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[sessionToken] = session
}

func (s *SessionsStore) Remove(sessionToken string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, sessionToken)
}

func (s *SessionsStore) Refresh(sessionToken string) (string, Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldSession := s.store[sessionToken]
	delete(s.store, sessionToken)
	newSession := Session{
		Username: oldSession.Username,
		IsAdmin:  oldSession.IsAdmin,
		Expiry:   time.Now().Add(480 * time.Second),
	}
	newToken := uuid.NewString()
	s.store[newToken] = newSession
	return newToken, newSession
}

func (s *SessionsStore) Get(sessionToken string) (Session, bool) {
	retrievedSession, ok := s.store[sessionToken]
	return retrievedSession, ok
}
