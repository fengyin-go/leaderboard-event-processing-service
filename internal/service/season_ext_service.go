package service

import (
	"time"

	"leaderboard/internal/model"
	"leaderboard/pkg/idgen"
)

func (s *Service) FinalizeSeason(id string) (*model.Season, error) {
	se, err := s.store.GetSeason(id)
	if err != nil {
		return nil, err
	}
	if se.Status == model.SeasonStatusFinished {
		return nil, model.NewValidationError("status", "赛季已结束")
	}
	if !model.SeasonCanTransition(se.Status, model.SeasonStatusFinished) {
		return nil, model.NewValidationError("status", "无法从当前状态结束赛季")
	}
	se.Status = model.SeasonStatusFinished
	se.EndAt = time.Now()
	if err := s.store.UpdateSeason(se); err != nil {
		return nil, err
	}
	_, _ = s.CaptureRankSnapshot(se.BoardID, se.ID)
	return se, nil
}

func (s *Service) CloneSeason(id string, newName string, startAt, endAt time.Time) (*model.Season, error) {
	if newName == "" {
		return nil, model.NewValidationError("name", "新赛季名称不能为空")
	}
	src, err := s.store.GetSeason(id)
	if err != nil {
		return nil, err
	}
	se := &model.Season{
		ID:        idgen.Hex(),
		BoardID:   src.BoardID,
		Name:      newName,
		StartAt:   startAt,
		EndAt:     endAt,
		Status:    model.SeasonStatusUpcoming,
		CreatedAt: time.Now(),
	}
	if err := se.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateSeason(se); err != nil {
		return nil, err
	}
	return se, nil
}

func (s *Service) GetCurrentSeason(boardID string) (*model.Season, error) {
	seasons := s.store.ListSeasons()
	for _, se := range seasons {
		if se.BoardID == boardID && se.Status == model.SeasonStatusActive {
			return se, nil
		}
	}
	return nil, model.NewValidationError("season", "当前没有活跃赛季")
}

func (s *Service) GetSeasonScoreEvents(seasonID string, page, size int) ([]*model.ScoreEvent, int, error) {
	return s.ListScoreEvents(model.ScoreEventFilter{SeasonID: seasonID}, page, size)
}

func (s *Service) GetSeasonRankEntries(seasonID string, page, size int) ([]*model.RankEntry, int, error) {
	return s.ListRankEntries(model.RankEntryFilter{SeasonID: seasonID}, page, size)
}
