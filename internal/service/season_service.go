package service

import (
	"sort"
	"time"

	"leaderboard/internal/model"
	"leaderboard/pkg/idgen"
)

func (s *Service) CreateSeason(input model.Season) (*model.Season, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetBoard(input.BoardID); err != nil {
		return nil, model.NewValidationError("board_id", "所属榜单不存在")
	}
	now := time.Now()
	se := &model.Season{
		ID:        idgen.Hex(),
		BoardID:   input.BoardID,
		Name:      input.Name,
		StartAt:   input.StartAt,
		EndAt:     input.EndAt,
		Status:    input.Status,
		CreatedAt: now,
	}
	if err := s.store.CreateSeason(se); err != nil {
		return nil, err
	}
	return se, nil
}

func (s *Service) ListSeasons(filter model.SeasonFilter, page, size int) ([]*model.Season, int, error) {
	all := s.store.ListSeasons()
	matched := make([]*model.Season, 0, len(all))
	for _, se := range all {
		if filter.Match(se) {
			matched = append(matched, se)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Season{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetSeason(id string) (*model.Season, error) {
	return s.store.GetSeason(id)
}

func (s *Service) UpdateSeason(id string, input model.Season) (*model.Season, error) {
	se, err := s.store.GetSeason(id)
	if err != nil {
		return nil, err
	}
	se.Name = input.Name
	se.StartAt = input.StartAt
	se.EndAt = input.EndAt
	if err := se.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSeason(se); err != nil {
		return nil, err
	}
	return se, nil
}

func (s *Service) TransitionSeasonStatus(id string, toStatus string) (*model.Season, error) {
	se, err := s.store.GetSeason(id)
	if err != nil {
		return nil, err
	}
	if !model.SeasonCanTransition(se.Status, toStatus) {
		return nil, model.NewValidationError("status", "状态流转不合法")
	}
	se.Status = toStatus
	if err := s.store.UpdateSeason(se); err != nil {
		return nil, err
	}
	return se, nil
}

func (s *Service) DeleteSeason(id string) error {
	return s.store.DeleteSeason(id)
}
