package service

import (
	"log/slog"
	"time"
)

func (c *Service) GetTrendingTopics() ([]string, error) {
	topics, err := c.db.GetTrendingTopics()

	if err != nil {
		return []string{}, err
	}

	slog.Info("Trending topics retrieved: ", "time", time.Now(), "count", len(topics), "topics", topics)

	return topics, nil
}