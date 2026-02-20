package models

import (
	"time"
)

// UserRole represents user role types
type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleHR       UserRole = "hr"
	UserRoleFinance  UserRole = "finance"
	UserRoleKaryawan UserRole = "karyawan"
)

// User represents admin user
type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Username   string    `json:"username" gorm:"unique;not null;size:50"`
	Password   string    `json:"-" gorm:"not null;size:255"`
	Name       string    `json:"name" gorm:"not null;size:100"`
	Email      string    `json:"email" gorm:"size:100"`
	Role       string    `json:"role" gorm:"default:'karyawan';size:20"`
	IsActive   bool      `json:"is_active" gorm:"default:true"`
	KaryawanID *uint     `json:"karyawan_id,omitempty"`
	Karyawan   *Karyawan `json:"karyawan,omitempty" gorm:"foreignKey:KaryawanID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// HasPermission checks if user has permission based on role hierarchy
func (u *User) HasPermission(requiredRole UserRole) bool {
	roleHierarchy := map[UserRole]int{
		UserRoleAdmin:    4,
		UserRoleFinance:  3,
		UserRoleHR:       2,
		UserRoleKaryawan: 1,
	}

	userLevel := roleHierarchy[UserRole(u.Role)]
	requiredLevel := roleHierarchy[requiredRole]

	return userLevel >= requiredLevel
}
