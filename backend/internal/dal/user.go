package dal

import (
	"errors"
	"ticktick-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserDAL 用户数据访问层
type UserDAL struct {
	db *Database
}

// NewUserDAL 创建用户数据访问层实例
func NewUserDAL(db *Database) *UserDAL {
	return &UserDAL{db: db}
}

// CreateUser 创建新用户
func (dal *UserDAL) CreateUser(user *models.User) error {
	return dal.db.GORM.Create(user).Error
}

// GetUserByEmail 根据邮箱获取用户
func (dal *UserDAL) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := dal.db.GORM.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID获取用户
func (dal *UserDAL) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := dal.db.GORM.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (dal *UserDAL) UpdateUser(user *models.User) error {
	return dal.db.GORM.Save(user).Error
}

// DeleteUser 软删除用户
func (dal *UserDAL) DeleteUser(id uuid.UUID) error {
	return dal.db.GORM.Delete(&models.User{}, id).Error
}

// EmailExists 检查邮箱是否已存在
func (dal *UserDAL) EmailExists(email string) (bool, error) {
	var count int64
	err := dal.db.GORM.Model(&models.User{}).Where("email = ? AND deleted_at IS NULL", email).Count(&count).Error
	return count > 0, err
}
