// ============================================
// user_service.go - 用户业务逻辑层
// ============================================
// 【知识点】Service 层是业务逻辑的核心
// 它调用 Repository 层获取/存储数据，编排业务流程
// Handler 层只负责 HTTP 请求/响应的处理，业务逻辑放在这里
//
// 为什么需要 Service 层?
// - 一个业务可能涉及多个 Repository（查用户 + 查角色 + 写缓存）
// - 多个 Handler 可能需要相同的业务逻辑（复用）
// - 业务逻辑变化时只改 Service，不影响 Handler

package service

import (
	"errors"
	"fmt"

	"go-vue-admin/server/internal/model"
	"go-vue-admin/server/internal/pkg"
	"go-vue-admin/server/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户业务逻辑
type UserService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
	jwtMgr   *pkg.JWTManager
}

// NewUserService 创建用户服务实例
// 【知识点】多个依赖注入: 需要用户仓库 + 角色仓库 + JWT管理器
func NewUserService(
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
	jwtMgr *pkg.JWTManager,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		jwtMgr:   jwtMgr,
	}
}

// RegisterInput 注册请求参数
type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
	Nickname string `json:"nickname"`
	// 【知识点】binding tag 是 Gin 的参数校验标签
	// required: 必填
	// min/max: 字符串长度范围
	// omitempty: 如果为空则跳过校验
	// email: 必须是合法邮箱格式
}

// LoginInput 登录请求参数
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
func (s *UserService) Register(input RegisterInput) (*model.User, error) {
	// 1. 检查用户名是否已存在
	_, err := s.userRepo.GetByUsername(input.Username)
	if err == nil {
		// 没有错误说明用户已存在
		return nil, errors.New("用户名已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 不是"记录不存在"的错误，说明数据库出问题了
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 2. 密码加密
	// 【知识点】bcrypt 是密码学推荐的哈希算法
	// - 自动加盐（salt），每次生成结果不同，防止彩虹表攻击
	// - 故意设计得很慢，增加暴力破解成本
	// - 永远不要用 MD5/SHA1 存密码！
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 3. 创建用户对象
	user := &model.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Email:    input.Email,
		Nickname: input.Nickname,
		Status:   1,
	}

	// 4. 分配默认角色
	// 【知识点】新注册用户默认给"普通用户"角色
	// admin 角色应该由管理员手动分配，不能自注册
	role, err := s.roleRepo.GetByCode("user")
	if err == nil {
		user.Roles = []model.Role{*role}
	}

	// 5. 保存到数据库
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return user, nil
}

// Login 用户登录
// 返回: 用户信息, JWT Token
func (s *UserService) Login(input LoginInput) (*model.User, string, error) {
	// 1. 查找用户
	user, err := s.userRepo.GetByUsername(input.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("用户名或密码错误")
		}
		return nil, "", fmt.Errorf("查询用户失败: %w", err)
	}

	// 2. 验证密码
	// 【知识点】CompareHashAndPassword 比较明文密码和哈希密码
	// 密码错误返回 bcrypt.ErrMismatchedHashAndPassword
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, "", errors.New("用户名或密码错误")
	}

	// 3. 检查用户状态
	if user.Status == 0 {
		return nil, "", errors.New("用户已被禁用")
	}

	// 4. 生成 JWT Token
	token, err := s.jwtMgr.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, "", fmt.Errorf("生成 Token 失败: %w", err)
	}

	return user, token, nil
}

// GetUser 获取单个用户详情
func (s *UserService) GetUser(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(query model.PageQuery) ([]model.User, int64, error) {
	return s.userRepo.List(query)
}

// UpdateUser 更新用户信息
// 【知识点】指针 vs 值: 传入 map[string]interface{} 可以只更新指定字段
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 如果更新了密码，需要重新加密
	if pwd, ok := updates["password"].(string); ok && pwd != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("密码加密失败: %w", err)
		}
		updates["password"] = string(hashedPassword)
	} else {
		delete(updates, "password") // 不允许通过 map 清空密码
	}

	// 【知识点】Gorm 的 Model().Updates() 只更新指定字段
	if err := s.userRepo.UpdateFields(user, updates); err != nil {
		return nil, err
	}

	return s.userRepo.GetByID(id)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}
