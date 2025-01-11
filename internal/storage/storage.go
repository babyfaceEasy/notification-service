package storage

import (
    "sync"
)

type NotificationStore struct {
    mu            sync.RWMutex
    notifications map[string][]Notification
}

func NewNotificationStore() *NotificationStore {
    return &NotificationStore{
        notifications: make(map[string][]Notification),
    }
}

func (s *NotificationStore) AddNotification(userID string, notification Notification) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.notifications[userID] = append(s.notifications[userID], notification)
}

func (s *NotificationStore) GetUserNotifications(userID string) []Notification {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.notifications[userID]
}

func (s *NotificationStore) ClearNotification(userID, notificationID string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    var filtered []Notification
    for _, n := range s.notifications[userID] {
        if n.ID != notificationID {
            filtered = append(filtered, n)
        }
    }
    s.notifications[userID] = filtered
}

func (s *NotificationStore) ClearAllNotifications(userID string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.notifications, userID)
}
