package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mariclezhang/vps_backend/internal/service"
	"github.com/mariclezhang/vps_backend/internal/util"
)

// NodeHandler 节点处理器
type NodeHandler struct {
	nodeService *service.NodeService
}

// NewNodeHandler 创建节点处理器实例
func NewNodeHandler() *NodeHandler {
	return &NodeHandler{
		nodeService: service.NewNodeService(),
	}
}

// GetList 获取节点列表
func (h *NodeHandler) GetList(c *gin.Context) {
	location := c.Query("location")
	protocol := c.Query("protocol")

	nodes, err := h.nodeService.GetNodes(location, protocol)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.Success(c, nodes)
}

// GetDetail 获取节点详情
func (h *NodeHandler) GetDetail(c *gin.Context) {
	idStr := c.Param("id")
	nodeID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.BadRequest(c, "无效的节点ID")
		return
	}

	node, err := h.nodeService.GetNodeDetail(nodeID)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.Success(c, node)
}

// TestLatency 测试节点延迟
func (h *NodeHandler) TestLatency(c *gin.Context) {
	idStr := c.Param("id")
	nodeID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.BadRequest(c, "无效的节点ID")
		return
	}

	latency, err := h.nodeService.TestLatency(nodeID)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.Success(c, gin.H{
		"latency": latency,
	})
}
