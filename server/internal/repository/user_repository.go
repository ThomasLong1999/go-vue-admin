// ============================================
// user_repository.go - 用户数据访问层
// ============================================
// 【知识点】Repository 层（仓库层）专门负责数据库操作
// 它把 SQL/Gorm 操作封装起来，Service 层不直接写数据库代码
// 好处:
// 1. Service 层不需要关心数据库细节（用什么数据库、怎么连接）
// 2. 更换数据库只需要改 Repository，不影响 Service
// 3. 方便写单元测试（可以 mock Repository）

package repository

import (
	"go-vue-admin/server/internal/model"

	"gorm.io/gorm"
)

// UserRepository 用户数据仓库
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
// 【知识点】依赖注入: 通过参数传入 db，而不是在内部创建
// 这样可以在测试时传入 mock 数据库
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
// 【知识点】&user 传入指针，Gorm 会把自增 ID 写回 user.ID
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据 ID 获取用户
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	// 【知识点】First 查询第一条匹配的记录
	// r.db.Where("id = ?", id).First(&user)
	// WHERE 中的 ? 是占位符，防止 SQL 注入
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户（登录时用）
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	// 【知识点】Preload 预加载关联数据
	// 不加 Preload 的话，user.Roles 会是 nil
	// 加了 Preload 会额外执行一条 SELECT 去查 user_roles 和 roles 表
	err := r.db.Preload("Roles").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表（分页）
// 【知识点】Gorm 的链式调用:
// db.Scope(...) → 应用作用域（复用查询条件）
// db.Offset(...).Limit(...) → 分页
// db.Find(&users) → 执行查询，结果写入 users
func (r *UserRepository) List(query model.PageQuery) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := r.db.Model(&model.User{})

	// 关键词搜索
	if query.Keyword != "" {
		// 【知识点】LIKE 模糊查询，% 匹配任意字符
		db = db.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	// 先查总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再查分页数据
	offset := query.GetOffset()
	pageSize := query.GetPageSize()
	err := db.Preload("Roles").Offset(offset).Limit(pageSize).
		Order("id DESC"). // 【知识点】按 ID 倒序，最新的排前面
		Find(&users).Error

	return users, total, err
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
	// 【知识点】Save 会更新所有字段（即使零值也会更新）
	// 如果只想更新非零值字段，用 Updates(user)
	// 但 Save 更安全，因为它不依赖零值判断
}

// Delete 删除用户（软删除，因为 User 有 DeletedAt 字段）
// 【知识点】软删除: 不执行 DELETE SQL，而是 UPDATE SET deleted_at = NOW()
// 之后查询时 Gorm 自动加 WHERE deleted_at IS NULL 过滤掉
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// UpdateFields 更新用户指定字段
// 只更新传入的非零值字段，密码需要单独处理
func (r *UserRepository) UpdateFields(user *model.User, updates map[string]interface{}) error {
	return r.db.Model(user).Updates(updates).Error
}

// GetByUsernameForAuth 获取用户密码（用于登录验证，不预加载角色）
func (r *UserRepository) GetByUsernameForAuth(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
