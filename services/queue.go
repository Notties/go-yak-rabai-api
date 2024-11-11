// services/queue.go
package services

import (
	"context"

	"yak.rabai/config"
	"yak.rabai/models"
)

// EnqueueUser adds a user to the queue based on their role (Speaker or Listener)
func EnqueueUser(userID string, role string) {
	key := "queue:" + role
	config.RDB.LPush(context.Background(), key, userID)
}

// MatchUsers checks if there is a Speaker and Listener available to create a room
func MatchUsers() *models.ChatRoom {
	speakerID, err := config.RDB.RPop(context.Background(), "queue:Speaker").Result()
	if err != nil {
		return nil
	}
	listenerID, err := config.RDB.RPop(context.Background(), "queue:Listener").Result()
	if err != nil {
		// If no listener is available, return speaker back to the queue
		EnqueueUser(speakerID, "Speaker")
		return nil
	}

	// Create a new ChatRoom
	chatRoom := models.ChatRoom{
		SpeakerID:  speakerID,
		ListenerID: listenerID,
	}
	config.DB.Create(&chatRoom)

	return &chatRoom
}

// GetQueueStatus returns the count of users waiting in each role's queue
func GetQueueStatus() map[string]int {
	speakerCount, _ := config.RDB.LLen(context.Background(), "queue:Speaker").Result()
	listenerCount, _ := config.RDB.LLen(context.Background(), "queue:Listener").Result()

	return map[string]int{
		"SpeakersWaiting":  int(speakerCount),
		"ListenersWaiting": int(listenerCount),
	}
}
