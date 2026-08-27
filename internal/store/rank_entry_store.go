package store

import (
	"leaderboard/internal/model"
)

func (s *MemoryStore) CreateRankEntry(r *model.RankEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rankEntries[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRankEntry(id string) (*model.RankEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rankEntries[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) GetRankEntryByBoardSeasonMember(boardID, seasonID, memberID string) (*model.RankEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.rankEntries {
		if r.BoardID == boardID && r.SeasonID == seasonID && r.MemberID == memberID {
			return r, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListRankEntries() []*model.RankEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RankEntry, 0, len(s.rankEntries))
	for _, r := range s.rankEntries {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateRankEntry(r *model.RankEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rankEntries[r.ID]; !ok {
		return ErrNotFound
	}
	s.rankEntries[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRankEntry(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rankEntries[id]; !ok {
		return ErrNotFound
	}
	delete(s.rankEntries, id)
	return nil
}
