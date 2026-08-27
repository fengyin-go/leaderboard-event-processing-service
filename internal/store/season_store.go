package store

import (
	"leaderboard/internal/model"
)

func (s *MemoryStore) CreateSeason(se *model.Season) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seasons[se.ID] = se
	return nil
}

func (s *MemoryStore) GetSeason(id string) (*model.Season, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	se, ok := s.seasons[id]
	if !ok {
		return nil, ErrNotFound
	}
	return se, nil
}

func (s *MemoryStore) ListSeasons() []*model.Season {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Season, 0, len(s.seasons))
	for _, se := range s.seasons {
		list = append(list, se)
	}
	return list
}

func (s *MemoryStore) UpdateSeason(se *model.Season) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.seasons[se.ID]; !ok {
		return ErrNotFound
	}
	s.seasons[se.ID] = se
	return nil
}

func (s *MemoryStore) DeleteSeason(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.seasons[id]; !ok {
		return ErrNotFound
	}
	delete(s.seasons, id)
	return nil
}
