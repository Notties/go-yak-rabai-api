// models/message.go
package models

import (
	"time"
)

type Message struct {
	ID         uint   `json:"id"`
	ChatRoomID string `json:"chat_room_id"` // Change this to string
	UserID     string `json:"user_id"`
	Content    string `json:"content"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// TableName specifies the table name for the Message model
func (Message) TableName() string {
	return "messages"
}
