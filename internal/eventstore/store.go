package eventstore

import (
	"reflect"
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

type Record019Handler interface{ Handle(string) string }
type Record019Processor struct{ prefix string }

func (p *Record019Processor) Handle(v string) string { return p.prefix + v }
func (s *EventStore) Record019(h Record019Handler, key string) string {
	// h == nil only catches an untyped nil interface. A typed-nil
	// interface (e.g. a nil *Record019Processor assigned to the
	// handler) is non-nil here yet still panics on Handle, which crashes
	// the processor and drops all subsequent score events. Treat both
	// forms of nil as a missing processor so callers get an explicit
	// result and processing continues.
	if h == nil || (reflect.ValueOf(h).Kind() == reflect.Ptr && reflect.ValueOf(h).IsNil()) {
		return "missing"
	}
	return h.Handle(key)
}
