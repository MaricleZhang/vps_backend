package model

import (
	"time"
)

// Announcement 公告模型
type Announcement struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Type      string    `json:"type" gorm:"default:'info'"` // info/warning/success
	Link      string    `json:"link"`
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Announcement) TableName() string {
	return "announcements"
}

// PasswordReset 密码重置模型
type PasswordReset struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"index;not null"`
	Code      string    `json:"code" gorm:"index;not null"`
	ExpiredAt time.Time `json:"expiredAt" gorm:"not null"`
	Used      bool      `json:"used" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// TableName 指定表名
func (PasswordReset) TableName() string {
	return "password_resets"
}
