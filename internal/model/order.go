package model

import (
	"time"
)

// Order 订单模型
type Order struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	UserID        int64     `json:"userId" gorm:"index"`
	OrderNo       string    `json:"orderNo" gorm:"uniqueIndex;not null"`
	Type          string    `json:"type"` // purchase/renew/recharge
	PlanID        *int64    `json:"planId"`
	Amount        float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	PaymentMethod string    `json:"paymentMethod"`
	Status        string    `json:"status" gorm:"default:'pending'"` // pending/paid/cancelled/refunded
	PaidAt        *time.Time `json:"paidAt"`
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	// Relations
	User *User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Plan *SubscriptionPlan `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}
