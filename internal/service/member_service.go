package service

import (
	"sort"
	"time"

	"leaderboard/internal/model"
	"leaderboard/pkg/idgen"
)

func (s *Service) CreateMember(input model.Member) (*model.Member, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	m := &model.Member{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Tag:       input.Tag,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateMember(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) ListMembers(filter model.MemberFilter, page, size int) ([]*model.Member, int, error) {
	all := s.store.ListMembers()
	matched := make([]*model.Member, 0, len(all))
	for _, m := range all {
		if filter.Match(m) {
			matched = append(matched, m)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Member{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetMember(id string) (*model.Member, error) {
	return s.store.GetMember(id)
}

func (s *Service) UpdateMember(id string, input model.Member) (*model.Member, error) {
	m, err := s.store.GetMember(id)
	if err != nil {
		return nil, err
	}
	m.Name = input.Name
	m.Tag = input.Tag
	m.UpdatedAt = time.Now()
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateMember(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) DeleteMember(id string) error {
	return s.store.DeleteMember(id)
}
