package rabbitmq

import (
	"GoEx21/internal/domain/model"
	"GoEx21/internal/reminder"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type companyEvent struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewCompanyEvent(channel *amqp.Channel) reminder.ReminderInterface {
	return &companyEvent{channel: channel}
}

func (c *companyEvent) Init(queueName string) error {
	var err error
	c.queue, err = c.channel.QueueDeclare(queueName, true, false, false, true, nil)
	if err != nil {
		return fmt.Errorf("rabbitmq.Init: %w", err)
	}
	return nil
}

func (c *companyEvent) SendEvent(company []model.Company) error {
	data, err := json.Marshal(company)
	if err != nil {
		return fmt.Errorf("rabbitmq.SendEvent: %w", err)
	}
	err = c.channel.Publish("", c.queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         data,
	})
	if err != nil {
		return fmt.Errorf("rabbitmq.SendEvent: %w", err)
	}
	return nil
}
