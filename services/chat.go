package services

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

var rooms = make(map[string]map[string]*websocket.Conn)
var mu sync.Mutex

// RegisterUserConnection stores user connections in rooms
func RegisterUserConnection(roomID, userID string, conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := rooms[roomID]; !exists {
		rooms[roomID] = make(map[string]*websocket.Conn)
	}
	rooms[roomID][userID] = conn
}

// ListenForMessages listens for incoming messages from a WebSocket connection
func ListenForMessages(conn *websocket.Conn, roomID, userID string) {
	defer conn.Close()

	for {
		// Read message from WebSocket
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		// Broadcast the message to all users in the room
		BroadcastMessage(roomID, userID, string(msg))
	}
}

// BroadcastMessage sends the message to all users in the room
func BroadcastMessage(roomID, userID, message string) {
	mu.Lock()
	defer mu.Unlock()

	for otherUserID, conn := range rooms[roomID] {
		if otherUserID != userID {
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				fmt.Println("Error sending message:", err)
			}
		}
	}
}
