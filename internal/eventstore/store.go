package eventstore

import "sync"

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

// Record028 records a ranking change for the given event key. The first
// submission records the change and returns true; any subsequent submission
// for the same key is a duplicate and returns false without appending another
// change, so re-reading the final ranking state is not duplicated.
//
// The done-check and the append are performed under a single lock hold so that
// concurrent duplicate submissions cannot both record a change.
func (s *EventStore) Record028(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.done[key] {
		return false
	}
	s.done[key] = true
	s.values[key] = append(s.values[key], "change")
	return true
}
