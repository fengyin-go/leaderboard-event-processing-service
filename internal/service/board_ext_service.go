package service

import (
	"sort"
	"time"

	"leaderboard/internal/model"
	"leaderboard/pkg/idgen"
)

func (s *Service) ArchiveBoard(id string) (*model.Board, error) {
	b, err := s.store.GetBoard(id)
	if err != nil {
		return nil, err
	}
	if b.Status == model.BoardStatusArchived {
		return nil, model.NewValidationError("status", "榜单已经是归档状态")
	}
	b.Status = model.BoardStatusArchived
	b.UpdatedAt = time.Now()
	if err := s.store.UpdateBoard(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) ActivateBoard(id string) (*model.Board, error) {
	b, err := s.store.GetBoard(id)
	if err != nil {
		return nil, err
	}
	if b.Status == model.BoardStatusActive {
		return nil, model.NewValidationError("status", "榜单已经是活跃状态")
	}
	b.Status = model.BoardStatusActive
	b.UpdatedAt = time.Now()
	if err := s.store.UpdateBoard(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) DuplicateBoard(id string, newName string) (*model.Board, error) {
	if newName == "" {
		return nil, model.NewValidationError("name", "新榜单名称不能为空")
	}
	src, err := s.store.GetBoard(id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	b := &model.Board{
		ID:          idgen.Hex(),
		Name:        newName,
		Description: src.Description,
		Metric:      src.Metric,
		SortOrder:   src.SortOrder,
		Status:      model.BoardStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateBoard(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) ListBoardsByMetric(metric string, page, size int) ([]*model.Board, int, error) {
	return s.ListBoards(model.BoardFilter{Metric: metric}, page, size)
}

func (s *Service) SearchBoards(keyword string, page, size int) ([]*model.Board, int, error) {
	return s.ListBoards(model.BoardFilter{Keyword: keyword}, page, size)
}

func (s *Service) GetBoardSeasons(boardID string) ([]*model.Season, error) {
	all := s.store.ListSeasons()
	var result []*model.Season
	for _, se := range all {
		if se.BoardID == boardID {
			result = append(result, se)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})
	return result, nil
}

func (s *Service) GetBoardStats(boardID string) (map[string]interface{}, error) {
	if _, err := s.store.GetBoard(boardID); err != nil {
		return nil, err
	}
	seasons, err := s.GetBoardSeasons(boardID)
	if err != nil {
		return nil, err
	}
	var upcoming, active, finished int
	for _, se := range seasons {
		switch se.Status {
		case model.SeasonStatusUpcoming:
			upcoming++
		case model.SeasonStatusActive:
			active++
		case model.SeasonStatusFinished:
			finished++
		}
	}
	entries := s.store.ListRankEntries()
	var memberCount int
	memberIDs := make(map[string]bool)
	for _, r := range entries {
		if r.BoardID == boardID {
			memberIDs[r.MemberID] = true
		}
	}
	memberCount = len(memberIDs)
	events := s.store.ListScoreEvents()
	var eventCount int
	for _, e := range events {
		if e.BoardID == boardID {
			eventCount++
		}
	}
	return map[string]interface{}{
		"board_id":       boardID,
		"season_count":   len(seasons),
		"upcoming":       upcoming,
		"active":         active,
		"finished":       finished,
		"member_count":   memberCount,
		"event_count":    eventCount,
	}, nil
}
