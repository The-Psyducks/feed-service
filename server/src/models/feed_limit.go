package models

import "strconv"


type LimitConfig struct {
	FromTime int
	Skip int
	Limit int
}

func NewLimitConfig(timeS string, skipS string, limitS string) LimitConfig {
	fromTime, _ := strconv.Atoi(timeS)
	skip, _ := strconv.Atoi(skipS)
	limit, _ := strconv.Atoi(limitS)

	return LimitConfig{
		FromTime: fromTime,
		Skip: skip,
		Limit: limit,
	}
}