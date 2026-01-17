package sessions

import (
	"sync"
	"time"
)

type Session struct {
	LastSeen time.Time
}

type Service struct {
	sessions map[string]*Session
	timeout  time.Duration
	mu       sync.RWMutex
}

func NewService(timeout time.Duration) *Service {
	return &Service{
		sessions: make(map[string]*Session),
		timeout:  timeout,
	}
}

func (s *Service) MarkActive(id string) {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[id]
	if ok {
		session.LastSeen = now
		return
	}
	s.sessions[id] = &Session{
		LastSeen: now,
	}
}

func (s *Service) IsActive(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[id]
	if !ok {
		return false
	}

	return time.Since(session.LastSeen) <= s.timeout
}

func (s *Service) CheckExpired() []string {
	now := time.Now()
	expired := []string{}

	s.mu.Lock()
	defer s.mu.Unlock()

	for id, sess := range s.sessions {
		if now.Sub(sess.LastSeen) >= s.timeout {
			expired = append(expired, id)
			delete(s.sessions, id)
		}
	}
	return expired
}
