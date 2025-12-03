package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/mariclezhang/vps_backend/internal/model"
	"github.com/mariclezhang/vps_backend/internal/util"
	"github.com/mariclezhang/vps_backend/pkg/db"
	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct{}

// NewAuthService 创建认证服务实例
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Login 用户登录
func (s *AuthService) Login(email, password string) (string, *model.User, error) {
	var user model.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("邮箱或密码错误")
		}
		return "", nil, err
	}

	// 检查密码
	if !util.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, errors.New("邮箱或密码错误")
	}

	// 检查用户状态
	if user.Status != "active" {
		return "", nil, errors.New("账户已被停用")
	}

	// 生成token
	token, err := util.GenerateToken(user.ID, user.Email, 24)
	if err != nil {
		return "", nil, err
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	db.DB.Save(&user)

	return token, &user, nil
}

// Register 用户注册
func (s *AuthService) Register(email, password, username string) error {
	// 检查邮箱是否已存在
	var count int64
	db.DB.Model(&model.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return errors.New("该邮箱已被注册")
	}

	// 加密密码
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return err
	}

	// 创建用户
	user := model.User{
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword,
		Balance:      0.00,
		Status:       "active",
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

// SendResetCode 发送重置密码验证码
func (s *AuthService) SendResetCode(email string) error {
	// 检查邮箱是否存在
	var user model.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("该邮箱未注册")
		}
		return err
	}

	// 生成6位验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 保存验证码（有效期15分钟）
	resetRecord := model.PasswordReset{
		Email:     email,
		Code:      code,
		ExpiredAt: time.Now().Add(15 * time.Minute),
		Used:      false,
	}

	if err := db.DB.Create(&resetRecord).Error; err != nil {
		return err
	}

	// TODO: 实际发送邮件
	// 这里仅打印日志，实际应该调用邮件服务
	fmt.Printf("重置密码验证码: %s (发送至 %s)\n", code, email)

	return nil
}

// ResetPassword 重置密码
func (s *AuthService) ResetPassword(email, code, newPassword string) error {
	// 查找有效的验证码
	var resetRecord model.PasswordReset
	if err := db.DB.Where("email = ? AND code = ? AND used = ? AND expired_at > ?",
		email, code, false, time.Now()).First(&resetRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("验证码无效或已过期")
		}
		return err
	}

	// 标记验证码为已使用
	resetRecord.Used = true
	db.DB.Save(&resetRecord)

	// 更新密码
	hashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := db.DB.Model(&model.User{}).Where("email = ?", email).
		Update("password_hash", hashedPassword).Error; err != nil {
		return err
	}

	return nil
}
