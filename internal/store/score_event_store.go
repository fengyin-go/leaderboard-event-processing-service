package store

import (
	"leaderboard/internal/model"
)

func (s *MemoryStore) CreateScoreEvent(e *model.ScoreEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scoreEvents[e.ID] = e
	return nil
}

func (s *MemoryStore) GetScoreEvent(id string) (*model.ScoreEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.scoreEvents[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *MemoryStore) ListScoreEvents() []*model.ScoreEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ScoreEvent, 0, len(s.scoreEvents))
	for _, e := range s.scoreEvents {
		list = append(list, e)
	}
	return list
}

func (s *MemoryStore) DeleteScoreEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.scoreEvents[id]; !ok {
		return ErrNotFound
	}
	delete(s.scoreEvents, id)
	return nil
}
