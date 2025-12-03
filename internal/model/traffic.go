package model

import (
	"time"
)

// TrafficLog 流量日志模型
type TrafficLog struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	UserID         int64     `json:"userId" gorm:"index:idx_traffic_user_recorded"`
	SubscriptionID int64     `json:"subscriptionId" gorm:"index"`
	NodeID         int64     `json:"nodeId"`
	UploadBytes    int64     `json:"uploadBytes" gorm:"default:0"`
	DownloadBytes  int64     `json:"downloadBytes" gorm:"default:0"`
	TotalBytes     int64     `json:"totalBytes" gorm:"default:0"`
	RecordedAt     time.Time `json:"recordedAt" gorm:"index:idx_traffic_user_recorded;autoCreateTime"`

	// Relations
	User         *User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Subscription *Subscription `json:"subscription,omitempty" gorm:"foreignKey:SubscriptionID"`
	Node         *Node         `json:"node,omitempty" gorm:"foreignKey:NodeID"`
}

// TableName 指定表名
func (TrafficLog) TableName() string {
	return "traffic_logs"
}
