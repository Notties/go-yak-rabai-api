// controllers/queue.go
package controllers

import (
	"strconv"

	"yak.rabai/services"

	"github.com/gin-gonic/gin"
)

// QueueChatRoom queues users based on their role and assigns rooms
func QueueChatRoom(c *gin.Context) {
	role := c.Query("role") // Role could be "Speaker" or "Listener"
	userID := c.Query("userID")

	if role == "Speaker" {
		services.EnqueueUser(userID, "Speaker")
	} else if role == "Listener" {
		services.EnqueueUser(userID, "Listener")
	}

	// Attempt to match users
	if match := services.MatchUsers(); match != nil {
		roomID := strconv.Itoa(int(match.ID))
		c.JSON(200, gin.H{"message": "Room found", "room_id": roomID})
	} else {
		c.JSON(200, gin.H{"message": "Waiting for match..."})
	}
}
