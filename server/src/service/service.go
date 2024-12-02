package service

import (
	"server/src/database"
)

const (
	FOLLOWING = "following"
	FORYOU    = "foryou"
	SINGLE    = "single"
	RETWEET   = "retweet"
)

type Service struct {
	db   database.Database
}

func NewService(db database.Database) *Service {
	return &Service{db: db}
}