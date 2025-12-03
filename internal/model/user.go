package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	Username     string    `json:"username" gorm:"not null"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;not null"`
	Avatar       string    `json:"avatar"`
	Balance      float64   `json:"balance" gorm:"type:decimal(10,2);default:0.00"`
	Status       string    `json:"status" gorm:"default:'active'"` // active/suspended/deleted
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	LastLoginAt  *time.Time `json:"lastLoginAt"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
