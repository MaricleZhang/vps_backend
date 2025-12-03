package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// StringArray 自定义类型用于存储字符串数组
type StringArray []string

// Scan 实现 sql.Scanner 接口
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}
	return json.Unmarshal(value.([]byte), a)
}

// Value 实现 driver.Valuer 接口
func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// SubscriptionPlan 订阅套餐模型
type SubscriptionPlan struct {
	ID           int64       `json:"id" gorm:"primaryKey"`
	Name         string      `json:"name" gorm:"not null"`
	Description  string      `json:"description" gorm:"type:text"`
	Price        float64     `json:"price" gorm:"type:decimal(10,2);not null"`
	TrafficLimit int64       `json:"traffic" gorm:"not null"` // 字节
	DurationDays int         `json:"duration" gorm:"not null"`
	Features     StringArray `json:"features" gorm:"type:jsonb"`
	IsActive     bool        `json:"isActive" gorm:"default:true"`
	SortOrder    int         `json:"sortOrder" gorm:"default:0"`
	CreatedAt    time.Time   `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time   `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}

// Subscription 用户订阅模型
type Subscription struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	UserID       int64     `json:"userId" gorm:"index;not null"`
	PlanID       int64     `json:"planId"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Status       string    `json:"status" gorm:"default:'active'"` // active/expired/cancelled
	TrafficLimit int64     `json:"traffic" gorm:"column:traffic_limit"`
	TrafficUsed  int64     `json:"trafficUsed" gorm:"default:0"`
	Price        float64   `json:"price" gorm:"type:decimal(10,2)"`
	DurationDays int       `json:"duration" gorm:"column:duration_days"`
	SubscribeURL string    `json:"subscribeUrl" gorm:"column:subscribe_url"`
	StartedAt    time.Time `json:"startedAt" gorm:"autoCreateTime"`
	ExpiredAt    time.Time `json:"expireDate"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Plan *SubscriptionPlan `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
}

// TableName 指定表名
func (Subscription) TableName() string {
	return "subscriptions"
}
