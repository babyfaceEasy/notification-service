package pubsub

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	var conn *amqp.Connection
	var err error

	// Retry logic for connecting to RabbitMQ
	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ connection failed: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Ensure the queue exists
	queueName := "notifications"
	_, err = ch.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}
	log.Printf("Queue '%s' is declared and ready.", queueName)

	return &RabbitMQ{conn: conn, channel: ch}, nil
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) Subscribe(queue string, handler func(msg []byte)) error {
	msgs, err := r.channel.Consume(
		queue,
		"",
		true,  // Auto-acknowledge
		false, // Not exclusive
		false, // No local
		false, // No wait
		nil,   // Arguments
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handler(msg.Body)
		}
	}()
	log.Printf("Subscribed to queue: %s", queue)
	return nil
}

func (r *RabbitMQ) Publish(queue, message string) error {
	_, err := r.channel.QueueDeclare(
		queue,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No wait
		nil,   // Arguments
	)
	if err != nil {
		return err
	}

	return r.channel.Publish(
		"",    // Exchange
		queue, // Routing key
		false, // Mandatory
		false, // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}
