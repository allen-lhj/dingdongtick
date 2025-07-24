package services

import (
	"errors"
	"fmt"
	"ticktick-backend/internal/dal"
	"ticktick-backend/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	userDAL *dal.UserDAL
}

// NewUserService 创建用户服务实例
func NewUserService(db *dal.Database) *UserService {
	return &UserService{
		userDAL: dal.NewUserDAL(db),
	}
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 用户响应结构（不包含敏感信息）
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}

// Register 用户注册
func (s *UserService) Register(req *RegisterRequest) (*UserResponse, error) {
	// 检查邮箱是否已存在
	exists, err := s.userDAL.EmailExists(req.Email)
	if err != nil {
		return nil, fmt.Errorf("检查邮箱失败: %w", err)
	}
	if exists {
		return nil, errors.New("邮箱已被注册")
	}

	// 哈希密码
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码哈希失败: %w", err)
	}

	// 创建用户
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	if err := s.userDAL.CreateUser(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return s.toUserResponse(user), nil
}

// Login 用户登录
func (s *UserService) Login(req *LoginRequest) (*UserResponse, error) {
	// 根据邮箱查找用户
	user, err := s.userDAL.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("查找用户失败: %w", err)
	}
	if user == nil {
		return nil, errors.New("邮箱或密码错误")
	}

	// 验证密码
	if !s.checkPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("邮箱或密码错误")
	}

	return s.toUserResponse(user), nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uuid.UUID) (*UserResponse, error) {
	user, err := s.userDAL.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("查找用户失败: %w", err)
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	return s.toUserResponse(user), nil
}

// hashPassword 哈希密码
func (s *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPassword 验证密码
func (s *UserService) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// toUserResponse 转换为用户响应结构
func (s *UserService) toUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
