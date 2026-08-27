package eventstore

import "sync"

// EventStore owns score event state shared by request processors.
type EventStore struct {
	mu      sync.RWMutex
	values  map[string][]string
	pending map[string][]string // staged, not-yet-committed events; invisible to Snapshot
	attempt map[string]int
	done    map[string]bool
}

func New() *EventStore {
	return &EventStore{
		values:  map[string][]string{},
		pending: map[string][]string{},
		attempt: map[string]int{},
		done:    map[string]bool{},
	}
}

func (s *EventStore) Append(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = append(s.values[key], value)
}

// Snapshot returns only the confirmed values for key. Pending events staged
// via StagePending are deliberately excluded so that a read performed right
// after a pending write never observes uncommitted intermediate values.
func (s *EventStore) Snapshot(key string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]string(nil), s.values[key]...)
}

// StagePending stages a not-yet-committed event without exposing it through
// Snapshot. It must not mutate any previously returned snapshot.
func (s *EventStore) StagePending(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pending[key] = append(s.pending[key], value)
}

// CommitPending merges staged pending events into the confirmed values and
// clears the pending staging area for the key.
func (s *EventStore) CommitPending(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = append(s.values[key], s.pending[key]...)
	delete(s.pending, key)
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

// Record022 stages a pending ("not-yet-committed") event for key without
// polluting confirmed reads. The returned view is the confirmed snapshot at
// call time and never reflects the pending event — exactly what a read right
// after a pending write should observe. The done channel is closed once the
// pending event has been staged; it is not committed unless CommitPending is
// called, so a subsequent Snapshot still returns only confirmed values.
func (s *EventStore) Record022(key string, ready <-chan struct{}) ([]string, <-chan struct{}) {
	view := s.Snapshot(key)
	done := make(chan struct{})
	go func() {
		<-ready
		s.StagePending(key, "pending-22")
		close(done)
	}()
	return view, done
}
