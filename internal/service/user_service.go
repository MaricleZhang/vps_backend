package service

import (
	"errors"

	"github.com/mariclezhang/vps_backend/internal/model"
	"github.com/mariclezhang/vps_backend/internal/util"
	"github.com/mariclezhang/vps_backend/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserService 用户服务
type UserService struct{}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{}
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID int64) (*model.User, error) {
	var user model.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUserInfo 更新用户信息
func (s *UserService) UpdateUserInfo(userID int64, username, avatar string) error {
	updates := make(map[string]interface{})

	if username != "" {
		updates["username"] = username
	}
	if avatar != "" {
		updates["avatar"] = avatar
	}

	if len(updates) == 0 {
		return errors.New("没有要更新的内容")
	}

	if err := db.DB.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID int64, oldPassword, newPassword string) error {
	var user model.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return err
	}

	// 验证旧密码
	if !util.CheckPasswordHash(oldPassword, user.PasswordHash) {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// 更新密码
	if err := db.DB.Model(&user).Update("password_hash", hashedPassword).Error; err != nil {
		return err
	}

	return nil
}

// GetBalance 获取用户余额
func (s *UserService) GetBalance(userID int64) (float64, error) {
	var user model.User
	if err := db.DB.Select("balance").First(&user, userID).Error; err != nil {
		return 0, err
	}
	return user.Balance, nil
}

// AddBalance 增加余额
func (s *UserService) AddBalance(userID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("金额必须大于0")
	}

	if err := db.DB.Model(&model.User{}).Where("id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
		return err
	}

	return nil
}

// DeductBalance 扣除余额
func (s *UserService) DeductBalance(userID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("金额必须大于0")
	}

	// 使用事务确保余额充足
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
			return err
		}

		if user.Balance < amount {
			return errors.New("余额不足")
		}

		if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		return nil
	})
}
