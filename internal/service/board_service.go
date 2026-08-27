package service

import (
	"sort"
	"time"

	"leaderboard/internal/model"
	"leaderboard/pkg/idgen"
)

func (s *Service) CreateBoard(input model.Board) (*model.Board, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	b := &model.Board{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Description: input.Description,
		Metric:      input.Metric,
		SortOrder:   input.SortOrder,
		Status:      input.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateBoard(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) ListBoards(filter model.BoardFilter, page, size int) ([]*model.Board, int, error) {
	all := s.store.ListBoards()
	matched := make([]*model.Board, 0, len(all))
	for _, b := range all {
		if filter.Match(b) {
			matched = append(matched, b)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Board{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetBoard(id string) (*model.Board, error) {
	return s.store.GetBoard(id)
}

func (s *Service) UpdateBoard(id string, input model.Board) (*model.Board, error) {
	b, err := s.store.GetBoard(id)
	if err != nil {
		return nil, err
	}
	b.Name = input.Name
	b.Description = input.Description
	b.Metric = input.Metric
	b.SortOrder = input.SortOrder
	b.Status = input.Status
	b.UpdatedAt = time.Now()
	if err := b.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateBoard(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) DeleteBoard(id string) error {
	return s.store.DeleteBoard(id)
}
