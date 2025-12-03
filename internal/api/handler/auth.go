package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mariclezhang/vps_backend/internal/middleware"
	"github.com/mariclezhang/vps_backend/internal/service"
	"github.com/mariclezhang/vps_backend/internal/util"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username"`
}

// SendResetCodeRequest 发送重置验证码请求
type SendResetCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	// 如果没有提供用户名，使用邮箱前缀
	if req.Username == "" {
		req.Username = req.Email[:len(req.Email)-len("@example.com")]
	}

	if err := h.authService.Register(req.Email, req.Password, req.Username); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "注册成功", gin.H{
		"success": true,
	})
}

// SendResetCode 发送重置密码验证码
func (h *AuthHandler) SendResetCode(c *gin.Context) {
	var req SendResetCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.authService.SendResetCode(req.Email); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "验证码已发送", gin.H{
		"success": true,
	})
}

// ResetPassword 重置密码
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.authService.ResetPassword(req.Email, req.Code, req.NewPassword); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "密码重置成功", gin.H{
		"success": true,
	})
}

// Logout 退出登录
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT是无状态的，客户端删除token即可
	// 如果需要token黑名单，可以在这里实现
	util.SuccessWithMessage(c, "退出成功", gin.H{
		"success": true,
	})
}

// GetCurrentUser 获取当前登录用户信息
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		util.Unauthorized(c, "未登录")
		return
	}

	util.Success(c, gin.H{
		"id": userID,
	})
}
