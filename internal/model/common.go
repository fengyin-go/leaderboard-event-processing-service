package model

import (
	"strings"
)

// TrimString 去除首尾空白，返回空字符串时也返回空。
func TrimString(s string) string {
	return strings.TrimSpace(s)
}

// IsEmptyString 判断字符串是否为空（去除空白后）。
func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}

// ClampInt 将整数限制在 [min, max] 范围内。
func ClampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// PageClamp 限制分页参数在合理范围。
func PageClamp(page, size, maxSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if maxSize > 0 && size > maxSize {
		size = maxSize
	}
	return page, size
}

// ListDiff 计算两个整数切片的差集（在 a 中但不在 b 中）。
func ListDiff(a, b []int64) []int64 {
	m := make(map[int64]bool, len(b))
	for _, v := range b {
		m[v] = true
	}
	var diff []int64
	for _, v := range a {
		if !m[v] {
			diff = append(diff, v)
		}
	}
	return diff
}

// ContainsString 判断字符串切片是否包含指定值。
func ContainsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// AllowedMetrics 返回合法的榜单指标列表。
func AllowedMetrics() []string {
	return []string{BoardMetricScore, BoardMetricStreak, BoardMetricCount}
}

// AllowedSortOrders 返回合法的排序方向列表。
func AllowedSortOrders() []string {
	return []string{BoardSortAsc, BoardSortDesc}
}

// AllowedBoardStatuses 返回合法的榜单状态列表。
func AllowedBoardStatuses() []string {
	return []string{BoardStatusActive, BoardStatusArchived}
}

// AllowedSeasonStatuses 返回合法的赛季状态列表。
func AllowedSeasonStatuses() []string {
	return []string{SeasonStatusUpcoming, SeasonStatusActive, SeasonStatusFinished}
}
