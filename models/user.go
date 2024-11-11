// models/user.go
package models

// User represents a user in the system
type User struct {
	ID        uint   `gorm:"primaryKey"`
	GoogleID  string `gorm:"uniqueIndex"`
	Name      string
	Email     string `gorm:"uniqueIndex"`
	CreatedAt int64
	UpdatedAt int64
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
