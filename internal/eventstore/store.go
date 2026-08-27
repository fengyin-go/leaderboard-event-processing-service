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

type Record026Handler interface{ Handle(string) string }
type Record026Processor struct{ prefix string }

func (p *Record026Processor) Handle(v string) string { return p.prefix + v }
func (s *EventStore) Record026(h Record026Handler, key string) string {
	// h 可能是 typed-nil interface：h == nil 检测不到，调用会解引用 nil 指针 panic。
	// 因此既判接口为空，也判其底层值为空，缺失处理器时返回明确状态而非崩溃，
	// 让上层把这条得分记为失败并继续处理后续事件。
	if h == nil || reflect.ValueOf(h).IsNil() {
		return "missing"
	}
	return h.Handle(key)
}
