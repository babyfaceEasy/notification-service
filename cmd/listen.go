package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

var (
	//rabbitMQURL string
	//queueName   string
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for notifications",
	Long:  "Listen for notifications on the RabbitMQ queue and log them to the console",
	Run: func(cmd *cobra.Command, args []string) {
		rabbitMQURL := os.Getenv("RABBITMQ_URL")
		queueName := os.Getenv("QUEUE_NAME")

		if rabbitMQURL == "" || queueName == "" {
			log.Fatalf("Environment variables RABBITMQ_URL or QUEUE_NAME are not set")
		}

		listenForNotifications(rabbitMQURL, queueName)
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)

	// Flags for RabbitMQ connection
	//listenCmd.Flags().StringVar(&rabbitMQURL, "rabbitmq-url", "amqp://guest:guest@localhost:5672/", "RabbitMQ URL")
	//listenCmd.Flags().StringVar(&queueName, "queue", "notifications", "Name of the RabbitMQ queue to listen to")
}

func listenForNotifications(rabbitMQURL, queueName string) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue (ensure it exists)
	_, err = ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Consume messages
	msgs, err := ch.Consume(
		queueName, // queue name
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Printf("Listening for notifications on queue: %s", queueName)

	// Handle messages
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received notification: %s", d.Body)
		}
	}()

	<-forever
}
