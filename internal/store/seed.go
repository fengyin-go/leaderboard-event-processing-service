package store

import (
	"time"

	"leaderboard/internal/model"
	"leaderboard/pkg/idgen"
)

func SeedDefaultData(s Store) error {
	boards := []*model.Board{
		{ID: idgen.Hex(), Name: "综合积分榜", Description: "综合表现排行", Metric: model.BoardMetricScore, SortOrder: model.BoardSortDesc, Status: model.BoardStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: idgen.Hex(), Name: "连胜榜", Description: "连续胜利场次", Metric: model.BoardMetricStreak, SortOrder: model.BoardSortDesc, Status: model.BoardStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: idgen.Hex(), Name: "参与次数榜", Description: "累计参与次数", Metric: model.BoardMetricCount, SortOrder: model.BoardSortDesc, Status: model.BoardStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, b := range boards {
		if _, err := s.GetBoardByName(b.Name); err == nil {
			continue
		}
		if err := s.CreateBoard(b); err != nil {
			return err
		}
	}

	members := []*model.Member{
		{ID: idgen.Hex(), Name: "Alice", Tag: "alice", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: idgen.Hex(), Name: "Bob", Tag: "bob", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: idgen.Hex(), Name: "Charlie", Tag: "charlie", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: idgen.Hex(), Name: "David", Tag: "david", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: idgen.Hex(), Name: "Eve", Tag: "eve", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, m := range members {
		if _, err := s.GetMemberByTag(m.Tag); err == nil {
			continue
		}
		if err := s.CreateMember(m); err != nil {
			return err
		}
	}

	boardList := s.ListBoards()
	memberList := s.ListMembers()
	if len(boardList) > 0 && len(memberList) > 0 {
		b := boardList[0]
		se := &model.Season{
			ID:        idgen.Hex(),
			BoardID:   b.ID,
			Name:      "2024 春季赛",
			StartAt:   time.Now().AddDate(0, 0, -30),
			EndAt:     time.Now().AddDate(0, 0, 60),
			Status:    model.SeasonStatusActive,
			CreatedAt: time.Now(),
		}
		if err := s.CreateSeason(se); err != nil {
			return err
		}
		for i, m := range memberList {
			e := &model.ScoreEvent{
				ID:       idgen.Hex(),
				BoardID:  b.ID,
				MemberID: m.ID,
				SeasonID: se.ID,
				Value:    int64(100 - i*15),
				Source:   "seed",
				At:       time.Now().AddDate(0, 0, -i),
			}
			if err := s.CreateScoreEvent(e); err != nil {
				return err
			}
		}
	}
	return nil
}
