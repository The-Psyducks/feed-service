package service

import (
	// "log"
	"server/src/database"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	FOLLOWING = "following"
	FORYOU    = "foryou"
	SINGLE    = "single"
	RETWEET   = "retweet"
)

type Service struct {
	db   database.Database
	amqp *amqp.Channel
}

func NewService(db database.Database, queue *amqp.Channel) *Service {
	return &Service{db: db}
}
