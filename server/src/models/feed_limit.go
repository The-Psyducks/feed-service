package models

import "strconv"


type LimitConfig struct {
	FromTime string
	Skip int
	Limit int
}

func NewLimitConfig(fromTime string, skipS string, limitS string) LimitConfig {
	skip, _ := strconv.Atoi(skipS)
	limit, _ := strconv.Atoi(limitS)

	return LimitConfig{
		FromTime: fromTime,
		Skip: skip,
		Limit: limit,
	}
}