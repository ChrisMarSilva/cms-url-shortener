package helpers

import (
	"log"

	"github.com/streadway/amqp"
)

type IRabbitMQ interface {
	SendMessage(queue, message string) error
	GetChannel() *amqp.Channel
}

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ(connectionString string) IRabbitMQ {
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatal(err)
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	return &RabbitMQ{Connection: connection, Channel: channel}
}

func (r *RabbitMQ) SendMessage(queue, message string) error {
	_, err := r.Channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	if err := r.Channel.Publish("", queue, false, false, amqp.Publishing{DeliveryMode: amqp.Persistent, ContentType: "text/plain", Body: []byte(message)}); err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) GetChannel() *amqp.Channel {
	return r.Channel
}
