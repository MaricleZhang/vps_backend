package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mariclezhang/vps_backend/internal/middleware"
	"github.com/mariclezhang/vps_backend/internal/service"
	"github.com/mariclezhang/vps_backend/internal/util"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService         *service.UserService
	subscriptionService *service.SubscriptionService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService:         service.NewUserService(),
		subscriptionService: service.NewSubscriptionService(),
	}
}

// UpdateUserInfoRequest 更新用户信息请求
type UpdateUserInfoRequest struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// GetInfo 获取用户信息
func (h *UserHandler) GetInfo(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	user, err := h.userService.GetUserInfo(userID)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.Success(c, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"avatar":    user.Avatar,
		"createdAt": user.CreatedAt,
	})
}

// UpdateInfo 更新用户信息
func (h *UserHandler) UpdateInfo(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.userService.UpdateUserInfo(userID, req.Username, req.Avatar); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "更新成功", gin.H{
		"success": true,
	})
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "密码修改成功", gin.H{
		"success": true,
	})
}

// GetBalance 获取账户余额
func (h *UserHandler) GetBalance(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	balance, err := h.userService.GetBalance(userID)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.Success(c, gin.H{
		"balance":  balance,
		"currency": "CNY",
	})
}

// GetTraffic 获取流量使用情况
func (h *UserHandler) GetTraffic(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	traffic, err := h.subscriptionService.GetUserTraffic(userID)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.Success(c, traffic)
}

// GetStats 获取账户统计
func (h *UserHandler) GetStats(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	balance, _ := h.userService.GetBalance(userID)
	traffic, _ := h.subscriptionService.GetUserTraffic(userID)

	util.Success(c, gin.H{
		"balance": balance,
		"traffic": traffic,
	})
}
