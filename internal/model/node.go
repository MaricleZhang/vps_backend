package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// NodeConfig 节点配置
type NodeConfig map[string]interface{}

// Scan 实现 sql.Scanner 接口
func (c *NodeConfig) Scan(value interface{}) error {
	if value == nil {
		*c = NodeConfig{}
		return nil
	}
	return json.Unmarshal(value.([]byte), c)
}

// Value 实现 driver.Valuer 接口
func (c NodeConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Node 节点模型
type Node struct {
	ID                 int64      `json:"id" gorm:"primaryKey"`
	Name               string     `json:"name" gorm:"not null"`
	Location           string     `json:"location"` // 地理位置
	Protocol           string     `json:"protocol"` // vmess/vless/trojan/shadowsocks
	ServerAddress      string     `json:"serverAddress" gorm:"column:server_address"`
	ServerPort         int        `json:"serverPort" gorm:"column:server_port"`
	Status             string     `json:"status" gorm:"default:'online'"` // online/offline/maintenance
	Latency            int        `json:"latency"`                        // 延迟 ms
	LoadPercentage     int        `json:"load" gorm:"column:load_percentage;default:0"`
	Bandwidth          string     `json:"bandwidth"`
	MaxConnections     int        `json:"maxConnections" gorm:"column:max_connections"`
	CurrentConnections int        `json:"currentConnections" gorm:"column:current_connections;default:0"`
	Config             NodeConfig `json:"config" gorm:"type:jsonb"`
	IsActive           bool       `json:"isActive" gorm:"default:true"`
	CreatedAt          time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt          time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Node) TableName() string {
	return "nodes"
}

// UserNodeAccess 用户节点访问权限模型
type UserNodeAccess struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	UserID         int64     `json:"userId" gorm:"uniqueIndex:idx_user_node"`
	NodeID         int64     `json:"nodeId" gorm:"uniqueIndex:idx_user_node"`
	SubscriptionID int64     `json:"subscriptionId"`
	GrantedAt      time.Time `json:"grantedAt" gorm:"autoCreateTime"`
	ExpiredAt      *time.Time `json:"expiredAt"`

	// Relations
	User         *User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Node         *Node         `json:"node,omitempty" gorm:"foreignKey:NodeID"`
	Subscription *Subscription `json:"subscription,omitempty" gorm:"foreignKey:SubscriptionID"`
}

// TableName 指定表名
func (UserNodeAccess) TableName() string {
	return "user_node_access"
}
