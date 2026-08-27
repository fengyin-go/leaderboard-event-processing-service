package eventstore

import (
	"context"
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

func (s *EventStore) Record024(key string, ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 每次只记录本次请求 ctx 的状态：取消只影响被取消的那次请求，
	// 下一条带全新 ctx 的正常得分事件会记录自己的（无错误）状态，
	// 不会把上一次已取消请求的状态带到下一条事件里。
	err := ""
	if ctx.Err() != nil {
		err = ctx.Err().Error()
	}
	s.values[key] = append(s.values[key], err)
}
