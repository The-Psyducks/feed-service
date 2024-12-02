package service

import (
	"server/src/models"
)

func (c *Service) GetUserMetrics(userID string, limits models.MetricLimits) (models.UserMetrics, error) {
	metrics, err := c.db.GetUserMetrics(userID, limits)

	if err != nil {
		return models.UserMetrics{}, err
	}

	return metrics, nil
}