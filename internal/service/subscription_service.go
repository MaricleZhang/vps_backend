package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/mariclezhang/vps_backend/internal/model"
	"github.com/mariclezhang/vps_backend/pkg/db"
	"gorm.io/gorm"
)

// SubscriptionService 订阅服务
type SubscriptionService struct {
	userService *UserService
}

// NewSubscriptionService 创建订阅服务实例
func NewSubscriptionService() *SubscriptionService {
	return &SubscriptionService{
		userService: NewUserService(),
	}
}

// GetUserSubscriptions 获取用户订阅列表
func (s *SubscriptionService) GetUserSubscriptions(userID int64) ([]model.Subscription, error) {
	var subscriptions []model.Subscription
	if err := db.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

// GetUserTraffic 获取用户流量使用情况
func (s *SubscriptionService) GetUserTraffic(userID int64) (map[string]interface{}, error) {
	// 获取用户的活跃订阅
	var subscription model.Subscription
	if err := db.DB.Where("user_id = ? AND status = ?", userID, "active").
		Order("expired_at DESC").
		First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{}{
				"used":       0,
				"total":      0,
				"percentage": 0,
				"resetDate":  nil,
			}, nil
		}
		return nil, err
	}

	percentage := 0.0
	if subscription.TrafficLimit > 0 {
		percentage = float64(subscription.TrafficUsed) / float64(subscription.TrafficLimit) * 100
	}

	return map[string]interface{}{
		"used":       subscription.TrafficUsed,
		"total":      subscription.TrafficLimit,
		"percentage": percentage,
		"resetDate":  subscription.ExpiredAt,
	}, nil
}

// GetAllPlans 获取所有可用套餐
func (s *SubscriptionService) GetAllPlans() ([]model.SubscriptionPlan, error) {
	var plans []model.SubscriptionPlan
	if err := db.DB.Where("is_active = ?", true).
		Order("sort_order ASC, price ASC").
		Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

// PurchaseSubscription 购买订阅
func (s *SubscriptionService) PurchaseSubscription(userID, planID int64, paymentMethod string) (*model.Subscription, error) {
	// 获取套餐信息
	var plan model.SubscriptionPlan
	if err := db.DB.First(&plan, planID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("套餐不存在")
		}
		return nil, err
	}

	if !plan.IsActive {
		return nil, errors.New("该套餐已下架")
	}

	// 开始事务
	var subscription *model.Subscription
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 扣除余额
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		if user.Balance < plan.Price {
			return errors.New("余额不足")
		}

		if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", plan.Price)).Error; err != nil {
			return err
		}

		// 创建订阅
		now := time.Now()
		expiredAt := now.AddDate(0, 0, plan.DurationDays)

		subscription = &model.Subscription{
			UserID:       userID,
			PlanID:       plan.ID,
			Name:         plan.Name,
			Type:         "monthly",
			Status:       "active",
			TrafficLimit: plan.TrafficLimit,
			TrafficUsed:  0,
			Price:        plan.Price,
			DurationDays: plan.DurationDays,
			SubscribeURL: s.generateSubscribeURL(userID),
			StartedAt:    now,
			ExpiredAt:    expiredAt,
		}

		if err := tx.Create(subscription).Error; err != nil {
			return err
		}

		// 创建订单记录
		orderNo := s.generateOrderNo(userID)
		order := model.Order{
			UserID:        userID,
			OrderNo:       orderNo,
			Type:          "purchase",
			PlanID:        &plan.ID,
			Amount:        plan.Price,
			PaymentMethod: paymentMethod,
			Status:        "paid",
			PaidAt:        &now,
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// 分配节点访问权限
		if err := s.grantNodeAccess(tx, userID, subscription.ID, expiredAt); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// RenewSubscription 续费订阅
func (s *SubscriptionService) RenewSubscription(userID, subscriptionID int64, months int) error {
	if months <= 0 {
		return errors.New("续费月数必须大于0")
	}

	// 获取订阅信息
	var subscription model.Subscription
	if err := db.DB.First(&subscription, subscriptionID).Error; err != nil {
		return err
	}

	if subscription.UserID != userID {
		return errors.New("无权操作此订阅")
	}

	// 获取套餐信息
	var plan model.SubscriptionPlan
	if err := db.DB.First(&plan, subscription.PlanID).Error; err != nil {
		return err
	}

	totalPrice := plan.Price * float64(months)

	// 开始事务
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// 扣除余额
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		if user.Balance < totalPrice {
			return errors.New("余额不足")
		}

		if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", totalPrice)).Error; err != nil {
			return err
		}

		// 延长过期时间
		newExpiredAt := subscription.ExpiredAt.AddDate(0, months, 0)
		if err := tx.Model(&subscription).Updates(map[string]interface{}{
			"expired_at": newExpiredAt,
			"status":     "active",
		}).Error; err != nil {
			return err
		}

		// 创建订单记录
		now := time.Now()
		orderNo := s.generateOrderNo(userID)
		order := model.Order{
			UserID:        userID,
			OrderNo:       orderNo,
			Type:          "renew",
			PlanID:        &plan.ID,
			Amount:        totalPrice,
			PaymentMethod: "balance",
			Status:        "paid",
			PaidAt:        &now,
		}

		return tx.Create(&order).Error
	})
}

// CancelSubscription 取消订阅
func (s *SubscriptionService) CancelSubscription(userID, subscriptionID int64) error {
	var subscription model.Subscription
	if err := db.DB.First(&subscription, subscriptionID).Error; err != nil {
		return err
	}

	if subscription.UserID != userID {
		return errors.New("无权操作此订阅")
	}

	if err := db.DB.Model(&subscription).Update("status", "cancelled").Error; err != nil {
		return err
	}

	return nil
}

// RecordTraffic 记录流量使用
func (s *SubscriptionService) RecordTraffic(userID, nodeID int64, uploadBytes, downloadBytes int64) error {
	totalBytes := uploadBytes + downloadBytes

	// 获取用户的活跃订阅
	var subscription model.Subscription
	if err := db.DB.Where("user_id = ? AND status = ?", userID, "active").
		Order("expired_at DESC").
		First(&subscription).Error; err != nil {
		return errors.New("没有活跃的订阅")
	}

	// 检查流量是否超限
	if subscription.TrafficUsed+totalBytes > subscription.TrafficLimit {
		return errors.New("流量已用完")
	}

	// 开始事务
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// 记录流量日志
		trafficLog := model.TrafficLog{
			UserID:         userID,
			SubscriptionID: subscription.ID,
			NodeID:         nodeID,
			UploadBytes:    uploadBytes,
			DownloadBytes:  downloadBytes,
			TotalBytes:     totalBytes,
		}

		if err := tx.Create(&trafficLog).Error; err != nil {
			return err
		}

		// 更新订阅的已用流量
		if err := tx.Model(&subscription).
			Update("traffic_used", gorm.Expr("traffic_used + ?", totalBytes)).Error; err != nil {
			return err
		}

		return nil
	})
}

// generateSubscribeURL 生成订阅链接
func (s *SubscriptionService) generateSubscribeURL(userID int64) string {
	// 使用用户ID和时间戳生成唯一token
	data := fmt.Sprintf("%d-%d", userID, time.Now().UnixNano())
	hash := md5.Sum([]byte(data))
	token := hex.EncodeToString(hash[:])
	return fmt.Sprintf("https://api.example.com/sub/%s", token)
}

// generateOrderNo 生成订单号
func (s *SubscriptionService) generateOrderNo(userID int64) string {
	return fmt.Sprintf("ORD%d%d", time.Now().Unix(), userID)
}

// grantNodeAccess 分配节点访问权限
func (s *SubscriptionService) grantNodeAccess(tx *gorm.DB, userID, subscriptionID int64, expiredAt time.Time) error {
	// 获取所有活跃节点
	var nodes []model.Node
	if err := tx.Where("is_active = ?", true).Find(&nodes).Error; err != nil {
		return err
	}

	// 为用户分配所有节点的访问权限
	for _, node := range nodes {
		access := model.UserNodeAccess{
			UserID:         userID,
			NodeID:         node.ID,
			SubscriptionID: subscriptionID,
			ExpiredAt:      &expiredAt,
		}

		// 使用 ON CONFLICT DO UPDATE 避免重复
		if err := tx.Create(&access).Error; err != nil {
			// 如果已存在，更新过期时间
			if err := tx.Model(&model.UserNodeAccess{}).
				Where("user_id = ? AND node_id = ?", userID, node.ID).
				Updates(map[string]interface{}{
					"subscription_id": subscriptionID,
					"expired_at":      &expiredAt,
				}).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
