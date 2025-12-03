package service

import (
	"testing"

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
	)
}

func TestAuthService_Register(t *testing.T) {
	setupTestDB(t)
	authService := NewAuthService()

	// 测试注册成功
	err := authService.Register("test@example.com", "password123", "testuser")
	assert.NoError(t, err)

	// 测试重复注册
	err = authService.Register("test@example.com", "password123", "testuser2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已被注册")
}

func TestAuthService_Login(t *testing.T) {
	setupTestDB(t)
	authService := NewAuthService()

	// 先注册用户
	authService.Register("test@example.com", "password123", "testuser")

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
