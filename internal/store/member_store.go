package store

import (
	"leaderboard/internal/model"
)

func (s *MemoryStore) CreateMember(m *model.Member) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.members {
		if exist.Tag == m.Tag {
			return ErrConflict
		}
	}
	s.members[m.ID] = m
	return nil
}

func (s *MemoryStore) GetMember(id string) (*model.Member, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.members[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

func (s *MemoryStore) GetMemberByTag(tag string) (*model.Member, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, m := range s.members {
		if m.Tag == tag {
			return m, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListMembers() []*model.Member {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Member, 0, len(s.members))
	for _, m := range s.members {
		list = append(list, m)
	}
	return list
}

func (s *MemoryStore) UpdateMember(m *model.Member) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.members[m.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.members {
		if exist.ID != m.ID && exist.Tag == m.Tag {
			return ErrConflict
		}
	}
	s.members[m.ID] = m
	return nil
}

func (s *MemoryStore) DeleteMember(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.members[id]; !ok {
		return ErrNotFound
	}
	delete(s.members, id)
	return nil
}
