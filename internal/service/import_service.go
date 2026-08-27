package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"leaderboard/internal/model"
)

type ImportResult struct {
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors,omitempty"`
}

func (s *Service) ImportMembersFromJSON(r io.Reader) (*ImportResult, error) {
	var inputs []model.Member
	dec := json.NewDecoder(r)
	if err := dec.Decode(&inputs); err != nil {
		return nil, model.NewValidationError("data", "JSON 解析失败: "+err.Error())
	}
	var result ImportResult
	for _, input := range inputs {
		_, err := s.CreateMember(input)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, err.Error())
		} else {
			result.Success++
		}
	}
	return &result, nil
}

func (s *Service) ImportScoreEventsFromJSON(r io.Reader) (*ImportResult, error) {
	var inputs []model.ScoreEvent
	dec := json.NewDecoder(r)
	if err := dec.Decode(&inputs); err != nil {
		return nil, model.NewValidationError("data", "JSON 解析失败: "+err.Error())
	}
	var result ImportResult
	for _, input := range inputs {
		_, err := s.CreateScoreEvent(input)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, err.Error())
		} else {
			result.Success++
		}
	}
	return &result, nil
}

func (s *Service) ImportMembersFromCSV(r io.Reader) (*ImportResult, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, model.NewValidationError("data", "CSV 解析失败: "+err.Error())
	}
	if len(records) < 1 {
		return nil, model.NewValidationError("data", "CSV 为空")
	}
	var result ImportResult
	for i, rec := range records {
		if i == 0 {
			continue
		}
		if len(rec) < 2 {
			result.Failed++
			result.Errors = append(result.Errors, "行 "+strconv.Itoa(i+1)+": 字段不足")
			continue
		}
		_, err := s.CreateMember(model.Member{Name: rec[0], Tag: rec[1]})
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, "行 "+strconv.Itoa(i+1)+": "+err.Error())
		} else {
			result.Success++
		}
	}
	return &result, nil
}

func (s *Service) ExportMembersToJSON() ([]byte, error) {
	members := s.store.ListMembers()
	return json.Marshal(members)
}

func (s *Service) ExportScoreEventsToJSON(boardID string) ([]byte, error) {
	all := s.store.ListScoreEvents()
	var filtered []*model.ScoreEvent
	for _, e := range all {
		if e.BoardID == boardID {
			filtered = append(filtered, e)
		}
	}
	return json.Marshal(filtered)
}

func (s *Service) ExportRankEntriesToCSV(boardID, seasonID string) ([]byte, error) {
	all := s.store.ListRankEntries()
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	_ = writer.Write([]string{"member_id", "score", "rank", "prev_rank"})
	for _, r := range all {
		if r.BoardID == boardID && r.SeasonID == seasonID {
			_ = writer.Write([]string{
				r.MemberID,
				strconv.FormatInt(r.Score, 10),
				strconv.Itoa(r.Rank),
				strconv.Itoa(r.PrevRank),
			})
		}
	}
	writer.Flush()
	return buf.Bytes(), nil
}

func (s *Service) ValidateImportData(data []byte) error {
	if len(data) == 0 {
		return errors.New("数据为空")
	}
	if len(data) > 10<<20 {
		return errors.New("数据超过 10MB 限制")
	}
	return nil
}
