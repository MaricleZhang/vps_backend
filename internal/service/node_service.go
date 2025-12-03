package service

import (
	"errors"

	"github.com/mariclezhang/vps_backend/internal/model"
	"github.com/mariclezhang/vps_backend/pkg/db"
	"gorm.io/gorm"
)

// NodeService 节点服务
type NodeService struct{}

// NewNodeService 创建节点服务实例
func NewNodeService() *NodeService {
	return &NodeService{}
}

// GetNodes 获取节点列表
func (s *NodeService) GetNodes(location, protocol string) ([]model.Node, error) {
	query := db.DB.Where("is_active = ?", true)

	if location != "" {
		query = query.Where("location = ?", location)
	}

	if protocol != "" {
		query = query.Where("protocol = ?", protocol)
	}

	var nodes []model.Node
	if err := query.Order("location ASC, name ASC").Find(&nodes).Error; err != nil {
		return nil, err
	}

	return nodes, nil
}

// GetNodeDetail 获取节点详情
func (s *NodeService) GetNodeDetail(nodeID int64) (*model.Node, error) {
	var node model.Node
	if err := db.DB.First(&node, nodeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("节点不存在")
		}
		return nil, err
	}

	return &node, nil
}

// TestLatency 测试节点延迟
func (s *NodeService) TestLatency(nodeID int64) (int, error) {
	var node model.Node
	if err := db.DB.First(&node, nodeID).Error; err != nil {
		return 0, err
	}

	// TODO: 实际的延迟测试逻辑
	// 这里仅模拟返回延迟值
	latency := 120 // ms

	// 更新节点延迟
	db.DB.Model(&node).Update("latency", latency)

	return latency, nil
}

// CheckUserNodeAccess 检查用户是否有权访问节点
func (s *NodeService) CheckUserNodeAccess(userID, nodeID int64) (bool, error) {
	var access model.UserNodeAccess
	err := db.DB.Where("user_id = ? AND node_id = ?", userID, nodeID).
		First(&access).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
