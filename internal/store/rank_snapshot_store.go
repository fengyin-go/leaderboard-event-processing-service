package store

import (
	"leaderboard/internal/model"
)

func (s *MemoryStore) CreateRankSnapshot(rs *model.RankSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rankSnapshots[rs.ID] = rs
	return nil
}

func (s *MemoryStore) GetRankSnapshot(id string) (*model.RankSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rs, ok := s.rankSnapshots[id]
	if !ok {
		return nil, ErrNotFound
	}
	return rs, nil
}

func (s *MemoryStore) ListRankSnapshots() []*model.RankSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RankSnapshot, 0, len(s.rankSnapshots))
	for _, rs := range s.rankSnapshots {
		list = append(list, rs)
	}
	return list
}

func (s *MemoryStore) DeleteRankSnapshot(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rankSnapshots[id]; !ok {
		return ErrNotFound
	}
	delete(s.rankSnapshots, id)
	return nil
}
