package service

func (c *Service) GetTrendingTopics() ([]string, error) {
	topics, err := c.db.GetTrendingTopics()

	if err != nil {
		return []string{}, err
	}

	return topics, nil
}