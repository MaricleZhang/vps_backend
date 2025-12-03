package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mariclezhang/vps_backend/internal/middleware"
	"github.com/mariclezhang/vps_backend/internal/service"
	"github.com/mariclezhang/vps_backend/internal/util"
)

// SubscriptionHandler 订阅处理器
type SubscriptionHandler struct {
	subscriptionService *service.SubscriptionService
}

// NewSubscriptionHandler 创建订阅处理器实例
func NewSubscriptionHandler() *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: service.NewSubscriptionService(),
	}
}

// PurchaseRequest 购买请求
type PurchaseRequest struct {
	PlanID        int64  `json:"planId" binding:"required"`
	PaymentMethod string `json:"paymentMethod" binding:"required"`
}

// RenewRequest 续费请求
type RenewRequest struct {
	SubscriptionID int64 `json:"subscriptionId" binding:"required"`
	Duration       int   `json:"duration" binding:"required,min=1"` // 月数
}

// GetList 获取用户订阅列表
func (h *SubscriptionHandler) GetList(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	subscriptions, err := h.subscriptionService.GetUserSubscriptions(userID)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.Success(c, subscriptions)
}

// GetPlans 获取可用套餐
func (h *SubscriptionHandler) GetPlans(c *gin.Context) {
	plans, err := h.subscriptionService.GetAllPlans()
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.Success(c, plans)
}

// Purchase 购买订阅
func (h *SubscriptionHandler) Purchase(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	subscription, err := h.subscriptionService.PurchaseSubscription(userID, req.PlanID, req.PaymentMethod)
	if err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "购买成功", subscription)
}

// Renew 续费订阅
func (h *SubscriptionHandler) Renew(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req RenewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.subscriptionService.RenewSubscription(userID, req.SubscriptionID, req.Duration); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "续费成功", gin.H{
		"success": true,
	})
}

// Cancel 取消订阅
func (h *SubscriptionHandler) Cancel(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	idStr := c.Param("id")
	subscriptionID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.BadRequest(c, "无效的订阅ID")
		return
	}

	if err := h.subscriptionService.CancelSubscription(userID, subscriptionID); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "取消成功", gin.H{
		"success": true,
	})
}

// Recharge 充值
func (h *SubscriptionHandler) Recharge(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
		Method string  `json:"method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "请求参数错误")
		return
	}

	// TODO: 实际支付流程
	// 这里仅模拟充值成功
	userService := service.NewUserService()
	if err := userService.AddBalance(userID, req.Amount); err != nil {
		util.Error(c, 400, err.Error())
		return
	}

	util.SuccessWithMessage(c, "充值成功", gin.H{
		"success": true,
		"amount":  req.Amount,
	})
}
