package sessions

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Session struct {
	LastSeen time.Time
}

type Service struct {
	sessions map[string]*Session
	timeout  time.Duration
	logger   zerolog.Logger
	mu       sync.RWMutex
}

type ServiceDeps struct {
	Timeout time.Duration
	Logger  zerolog.Logger
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		sessions: make(map[string]*Session),
		timeout:  deps.Timeout,
		logger:   deps.Logger,
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
	if len(expired) > 0 {
		s.logger.Info().Int("count", len(expired)).Msg("expired sessions cleaned")
	}
	return expired
}
