package service

import (
	"testing"
	"time"

	"github.com/mariclezhang/vps_backend/internal/model"
	"github.com/mariclezhang/vps_backend/pkg/db"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	var err error
	db.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 自动迁移
	db.DB.AutoMigrate(
		&model.User{},
		&model.SubscriptionPlan{},
		&model.Subscription{},
		&model.Node{},
		&model.UserNodeAccess{},
		&model.TrafficLog{},
		&model.Order{},
		&model.PasswordReset{},
	)
}

// createTestVerificationCode 创建测试验证码
func createTestVerificationCode(email, code string) {
	resetRecord := model.PasswordReset{
		Email:     email,
		Code:      code,
		ExpiredAt: time.Now().Add(15 * time.Minute),
		Used:      false,
	}
	db.DB.Create(&resetRecord)
}

func TestAuthService_SendRegisterCode(t *testing.T) {
	setupTestDB(t)
	authService := NewAuthService()

	// 测试发送注册验证码成功
	err := authService.SendRegisterCode("929006968@qq.com")
	assert.NoError(t, err)

	// 验证验证码已保存到数据库
	var resetRecord model.PasswordReset
	err = db.DB.Where("email = ?", "929006968@qq.com").First(&resetRecord).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, resetRecord.Code)
	assert.False(t, resetRecord.Used)

	// 测试已注册邮箱发送验证码失败
	user := model.User{
		Email:        "existing@example.com",
		Username:     "existing",
		PasswordHash: "hash",
		Status:       "active",
	}
	db.DB.Create(&user)

	err = authService.SendRegisterCode("existing@example.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已被注册")
}

func TestAuthService_Register(t *testing.T) {
	setupTestDB(t)
	authService := NewAuthService()

	// 创建测试验证码
	createTestVerificationCode("test@example.com", "123456")

	// 测试注册成功
	err := authService.Register("test@example.com", "password123", "123456")
	assert.NoError(t, err)

	// 验证用户已创建
	var user model.User
	err = db.DB.Where("email = ?", "test@example.com").First(&user).Error
	assert.NoError(t, err)
	assert.Equal(t, "test", user.Username) // 用户名应该是邮箱前缀

	// 验证验证码已标记为已使用
	var resetRecord model.PasswordReset
	db.DB.Where("email = ? AND code = ?", "test@example.com", "123456").First(&resetRecord)
	assert.True(t, resetRecord.Used)

	// 测试重复注册（使用新验证码）
	createTestVerificationCode("test@example.com", "654321")
	err = authService.Register("test@example.com", "password123", "654321")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已被注册")

	// 测试验证码无效
	err = authService.Register("newuser@example.com", "password123", "wrongcode")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "验证码无效或已过期")

	// 测试验证码已过期
	expiredRecord := model.PasswordReset{
		Email:     "expired@example.com",
		Code:      "111111",
		ExpiredAt: time.Now().Add(-1 * time.Hour), // 已过期
		Used:      false,
	}
	db.DB.Create(&expiredRecord)

	err = authService.Register("expired@example.com", "password123", "111111")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "验证码无效或已过期")
}

func TestAuthService_Login(t *testing.T) {
	setupTestDB(t)
	authService := NewAuthService()

	// 先创建验证码并注册用户
	createTestVerificationCode("test@example.com", "123456")
	authService.Register("test@example.com", "password123", "123456")

	// 测试登录成功
	token, user, err := authService.Login("test@example.com", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, "test@example.com", user.Email)

	// 测试密码错误
	_, _, err = authService.Login("test@example.com", "wrongpassword")
	assert.Error(t, err)

	// 测试用户不存在
	_, _, err = authService.Login("nonexistent@example.com", "password123")
	assert.Error(t, err)
}

func TestUserService_GetUserInfo(t *testing.T) {
	setupTestDB(t)
	userService := NewUserService()

	// 创建测试用户
	user := model.User{
		Email:    "test@example.com",
		Username: "testuser",
		Balance:  100.00,
		Status:   "active",
	}
	db.DB.Create(&user)

	// 测试获取用户信息
	result, err := userService.GetUserInfo(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, 100.00, result.Balance)

	// 测试用户不存在
	_, err = userService.GetUserInfo(999)
	assert.Error(t, err)
}

func TestUserService_Balance(t *testing.T) {
	setupTestDB(t)
	userService := NewUserService()

	// 创建测试用户
	user := model.User{
		Email:    "test@example.com",
		Username: "testuser",
		Balance:  100.00,
		Status:   "active",
	}
	db.DB.Create(&user)

	// 测试增加余额
	err := userService.AddBalance(user.ID, 50.00)
	assert.NoError(t, err)

	balance, _ := userService.GetBalance(user.ID)
	assert.Equal(t, 150.00, balance)

	// 测试扣除余额
	err = userService.DeductBalance(user.ID, 30.00)
	assert.NoError(t, err)

	balance, _ = userService.GetBalance(user.ID)
	assert.Equal(t, 120.00, balance)

	// 测试余额不足
	err = userService.DeductBalance(user.ID, 200.00)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "余额不足")
}
