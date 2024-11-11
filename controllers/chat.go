// controllers/chat.go
package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/websocket"
	"yak.rabai/config"
	"yak.rabai/lib/sockets"
	"yak.rabai/models"
	"yak.rabai/services"
)

var chatServer = sockets.ChatSocket()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ConnectChat serves the WebSocket connection
func ConnectChat(c *gin.Context) {
	roomID := c.DefaultQuery("room_id", "")
	userID := c.DefaultQuery("user_id", "")

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}

	// Store the connection in the room for broadcasting
	services.RegisterUserConnection(roomID, userID, conn)

	// Listen for incoming messages from the user
	go services.ListenForMessages(conn, roomID, userID)
}

// HandleChatMessage processes messages sent within a chat room
func HandleChatMessage(c *gin.Context) {
	var messageInput struct {
		RoomID  string `json:"room_id" binding:"required"`
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&messageInput); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Broadcast message to the chat room via Socket.IO
	chatServer.BroadcastToRoom("/", messageInput.RoomID, "message", messageInput.Message)

	// Save the message to the database
	userID := c.GetString("userID")
	message := models.Message{
		ChatRoomID: messageInput.RoomID,
		UserID:     userID,
		Content:    messageInput.Message,
	}

	if err := config.DB.Create(&message).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save message"})
		return
	}

	c.JSON(200, gin.H{"message": "Message sent successfully"})
}

// Handle typing indicator in a chat room
func HandleTypingIndicator(c *gin.Context) {
	var input struct {
		RoomID string `json:"room_id" binding:"required"`
		Typing bool   `json:"typing" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Broadcast typing status to the chat room
	chatServer.BroadcastToRoom("/", input.RoomID, "typing", input.Typing)

	c.JSON(200, gin.H{"message": "Typing status updated"})
}

// HandleReconnect manages reconnection to the same chat room
func HandleReconnect(s socketio.Conn, roomID string) {
	s.Join(roomID)
	chatServer.BroadcastToRoom("/", roomID, "reconnect", "User reconnected")

	// Send chat history
	var messages []models.Message
	config.DB.Where("chat_room_id = ?", roomID).Find(&messages)
	s.Emit("chat_history", messages)
}

// Handle user leaving the room and rating the listener
func LeaveRoomAndRate(c *gin.Context) {
	var input struct {
		RoomID  string `json:"room_id" binding:"required"`
		UserID  string `json:"user_id" binding:"required"`
		Rating  int    `json:"rating" binding:"required"`
		Comment string `json:"comment"`
	}

	// Bind the request data to the struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Save the rating
	if err := services.SaveRating(input.UserID, input.Rating, input.Comment); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save rating"})
		return
	}

	// Broadcast to the room that the user has left
	services.BroadcastMessage(input.RoomID, input.UserID, "User has left Chat")

	// Close the WebSocket connection
	chatServer.Close()

	// Respond to the client with success message
	c.JSON(200, gin.H{"message": "User left the room and rating saved successfully"})
}
