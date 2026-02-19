package models

import (
	"time"
)

// User represents admin user
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null;size:50"`
	Password  string    `json:"-" gorm:"not null;size:255"`
	Name      string    `json:"name" gorm:"not null;size:100"`
	Email     string    `json:"email" gorm:"size:100"`
	Role      string    `json:"role" gorm:"default:'admin';size:20"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
