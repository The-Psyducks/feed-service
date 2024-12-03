package models

import (
	"log/slog"
	"strconv"
)

const (
	MAX_LIMIT = 30
)


type LimitConfig struct {
	FromTime string
	Skip int
	Limit int
}

type MetricLimits struct {
	FromTime string
	ToTime string
}

func NewLimitConfig(fromTime string, skipS string, limitS string) LimitConfig {
	skip, _ := strconv.Atoi(skipS)
	limit, _ := strconv.Atoi(limitS)

	if limit > MAX_LIMIT {
		slog.Warn("Limit exceeds max limit", "limit", limit, "max_limit", MAX_LIMIT)
		limit = MAX_LIMIT
	}

	return LimitConfig{
		FromTime: fromTime,
		Skip: skip,
		Limit: limit,
	}
}
