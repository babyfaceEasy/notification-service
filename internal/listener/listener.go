package listener

import (
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/babyfaceeasy/notification_svc/internal/pubsub"
	"github.com/babyfaceeasy/notification_svc/internal/storage"
)

// StartListener starts the RabbitMQ listener for the notifications queue.
func StartListener(rabbitMQ *pubsub.RabbitMQ, store *storage.NotificationStore) error {
	// Subscribe to the "notifications" queue
	err := rabbitMQ.Subscribe("notifications", func(msg []byte) {
		log.Printf("Received message: %s", string(msg))
		notification := storage.Notification{
			ID:      generateID(),
			UserID:  "user1", // Extract from the message if needed
			Message: string(msg),
			Read:    false,
		}
		store.AddNotification(notification.UserID, notification)
	})
	if err != nil {
		return err
	}
	log.Println("Listener started successfully")
	return nil
}

// generateID generates a unique identifier for notifications.
// You can replace this with a UUID generator.
func generateID() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().String())))
	// return "unique-id" // Replace with UUID logic
}
