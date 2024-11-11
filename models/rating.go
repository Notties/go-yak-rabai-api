// models/rating.go
package models

// Rating represents the rating and comment given by the Speaker to the Listener
type Rating struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    string `gorm:"not null"`
	Rating    int    `gorm:"not null"`
	Comment   string `gorm:"not null"`
	CreatedAt int64
	UpdatedAt int64
}

// TableName specifies the table name for the Rating model
func (Rating) TableName() string {
	return "ratings"
}
