// models/chatroom.go
package models

// ChatRoom represents a room where users chat
type ChatRoom struct {
	ID         uint   `gorm:"primaryKey"`
	SpeakerID  string `gorm:"not null"`
	ListenerID string `gorm:"not null"`
	CreatedAt  int64
	UpdatedAt  int64
}

// TableName specifies the table name for the ChatRoom model
func (ChatRoom) TableName() string {
	return "chat_rooms"
}
