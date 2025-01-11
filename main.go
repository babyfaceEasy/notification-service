package main

import (
	"fmt"
	"log"

	"github.com/babyfaceeasy/notification_svc/cmd"
	"github.com/babyfaceeasy/notification_svc/internal/listener"
	"github.com/babyfaceeasy/notification_svc/internal/pubsub"
	"github.com/babyfaceeasy/notification_svc/internal/storage"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func mainOLd() {
	// Initialize storage
	store := storage.NewNotificationStore()

	// Initialize RabbitMQ
	rabbitMQ, err := pubsub.NewRabbitMQ("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Start the listener
	err = listener.StartListener(rabbitMQ, store)
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}

	// Start the REST server
	r := gin.Default()

	r.GET("/notifications/:user_id", func(c *gin.Context) {
		userID := c.Param("user_id")
		notifications := store.GetUserNotifications(userID)
		c.JSON(200, notifications)
	})

	r.DELETE("/notifications/:user_id/:notification_id", func(c *gin.Context) {
		userID := c.Param("user_id")
		notificationID := c.Param("notification_id")
		store.ClearNotification(userID, notificationID)
		c.Status(200)
	})

	r.DELETE("/notifications/:user_id", func(c *gin.Context) {
		userID := c.Param("user_id")
		store.ClearAllNotifications(userID)
		c.Status(200)
	})

	r.Run(":8080")
}

func main() {
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}
}

func loadEnv() error {
	return godotenv.Load(".env")
}

func initialize() error {
	if err := loadEnv(); err != nil {
		return fmt.Errorf("error loading environment variables: %v", err)
	}

	return nil
}
