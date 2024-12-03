package service

import (
	"log/slog"
	"server/src/models"
)

func (c *Service) GetUserMetrics(userID string, limits models.MetricLimits) (models.UserMetrics, error) {
	metrics, err := c.db.GetUserMetrics(userID, limits)

	if err != nil {
		return models.UserMetrics{}, err
	}

	slog.Info("User metrics retrieved: ", "user_id", userID)

	return metrics, nil
}