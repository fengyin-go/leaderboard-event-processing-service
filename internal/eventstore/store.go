package eventstore

import (
	"errors"
	"fmt"
	"sync"
)

// EventStore owns score event state shared by request processors.
type EventStore struct {
	mu      sync.RWMutex
	values  map[string][]string
	attempt map[string]int
	done    map[string]bool
}

func New() *EventStore {
	return &EventStore{values: map[string][]string{}, attempt: map[string]int{}, done: map[string]bool{}}
}

func (s *EventStore) Append(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = append(s.values[key], value)
}

func (s *EventStore) Snapshot(key string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]string(nil), s.values[key]...)
}

func (s *EventStore) SetAttempt(key string, n int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.attempt[key] = n
}

func (s *EventStore) Attempt(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.attempt[key]
}

func (s *EventStore) MarkDone(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.done[key] = true
}

func (s *EventStore) IsDone(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.done[key]
}

// Rollback drops the most recently appended value for key, undoing a write that
// failed as part of an isolated batch. It is a no-op when key has no values, so
// a rollback always balances a single preceding Append.
func (s *EventStore) Rollback(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v := s.values[key]
	if len(v) == 0 {
		return
	}
	s.values[key] = v[:len(v)-1]
}

func (s *EventStore) Record025(key string, attempt int) error {
	s.SetAttempt(key, attempt)
	s.Append(key, fmt.Sprint(attempt))
	if attempt == 1 {
		// The write was rejected with a temporary error, so roll it back to
		// keep the batch isolated. The retried attempt then produces the only
		// committed value, leaving a single result for the event.
		s.Rollback(key)
		return errors.New("temporary")
	}
	return nil
}
