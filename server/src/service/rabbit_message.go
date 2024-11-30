package service

import (
	"encoding/json"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"

	"server/src/models"
)

func sendMessage(queue *amqp.Channel, msg []byte) error {
	message := amqp.Publishing{
		ContentType:  "application/json",
		Body:         msg,
		DeliveryMode: amqp.Persistent,
	}

	fmt.Println("Sending login attempt to queue: ", message)
	err := queue.Publish("", os.Getenv("CLOUDAMQP_QUEUE"), false, false, message)

	if err != nil {
		return fmt.Errorf("error publishing login attempt message to queue: %w", err)
	}
	return nil
}

func (c *Service) sendNewContentMessage(hashtags []string, timestamp string) error {
	queueMsg := models.QueueMessage{
		MessageType: models.NEW_CONTENT,
		Message: models.New_Content{
			Hashtags:  hashtags,
			Timestamp: timestamp,
		},
	}

	newRegistry, err := json.Marshal(queueMsg)
	if err != nil {
		return fmt.Errorf("error marshalling login attempt message for rabbit: %w", err)
	}

	return sendMessage(c.amqp, newRegistry)
}
