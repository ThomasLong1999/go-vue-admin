// ============================================
// model.go - 数据模型定义
// ============================================
// 【知识点】Gorm 模型是用 struct 定义的，每个结构体对应数据库中的一张表
// 表名自动转换为蛇形命名（snake_case）:
// User → users, UserRole → user_roles
//
// Gorm 通过 struct tag 来定义字段的数据库属性:
// - gorm:"column:name"      → 指定列名
// - gorm:"primaryKey"       → 主键
// - gorm:"autoIncrement"    → 自增
// - gorm:"not null"         → 非空约束
// - gorm:"unique"           → 唯一约束
// - gorm:"size:255"         → 字符串长度
// - gorm:"type:varchar(32)" → 指定列类型
// - gorm:"index"            → 添加索引

package model

import (
	"time"

	"gorm.io/gorm"
)

// ============================================
// User 用户模型
// ============================================
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `json:"created_at"`                 // Gorm 自动管理的创建时间
	UpdatedAt time.Time      `json:"updated_at"`                 // Gorm 自动管理的更新时间
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`               // 【知识点】软删除字段
	// json:"-" 表示序列化时忽略这个字段，不返回给前端
	// gorm.DeletedAt 是一个可空时间，删除时只设置这个值，不真正删除记录
	// 查询时 Gorm 自动过滤掉已软删除的记录

	Username string `json:"username" gorm:"type:varchar(32);uniqueIndex;not null"`
	Password string `json:"-" gorm:"type:varchar(255);not null"`
	// json:"-" 不返回密码给前端！这是安全基础

	Email    string `json:"email" gorm:"type:varchar(128)"`
	Nickname string `json:"nickname" gorm:"type:varchar(64)"`
	Avatar   string `json:"avatar" gorm:"type:varchar(512)"`

	Status int8 `json:"status" gorm:"type:tinyint;default:1;comment:状态 1启用 0禁用"`
	// 【知识点】comment 是 MySQL 的列注释，方便在数据库工具中查看含义

	// 【知识点】Gorm 关联关系
	// many2many 表示多对多关系，自动创建中间表
	// 这里 User 和 Role 是多对多：一个用户可以有多个角色，一个角色可以分配给多个用户
	Roles []Role `json:"roles" gorm:"many2many:user_roles;"`
}

// ============================================
// Role 角色模型
// ============================================
// 【知识点】RBAC（Role-Based Access Control）基于角色的访问控制
// 用户 → 角色 → 权限，三层结构
// 例如: 用户"张三"有角色"管理员"，角色"管理员"有权限"删除用户"
type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Name        string         `json:"name" gorm:"type:varchar(32);uniqueIndex;not null"`
	Code        string         `json:"code" gorm:"type:varchar(32);uniqueIndex;not null"`
	// Code 是角色的标识码，如 "admin", "editor"
	Description string        `json:"description" gorm:"type:varchar(256)"`
	Status      int8           `json:"status" gorm:"type:tinyint;default:1"`

	// 角色和权限也是多对多关系
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}

// ============================================
// Permission 权限模型
// ============================================
type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Name        string         `json:"name" gorm:"type:varchar(64);not null"`
	Code        string         `json:"code" gorm:"type:varchar(64);uniqueIndex;not null"`
	// Code 格式: "user:list", "user:create", "user:delete" 等
	Description string         `json:"description" gorm:"type:varchar(256)"`
	Type        int8           `json:"type" gorm:"type:tinyint;default:1;comment:1菜单 2按钮"`
	// Type: 1=菜单权限（控制页面访问），2=按钮权限（控制操作按钮）
	ParentID    uint           `json:"parent_id" gorm:"default:0;comment:父权限ID"`
	// 【知识点】ParentID 实现树形结构
	// 权限可以分级: "用户管理"(父) → "查看列表"(子) / "新增"(子) / "删除"(子)
}

// ============================================
// PageQuery 分页查询参数（通用）
// ============================================
// 【知识点】Go 的"泛型"用方括号: type PageResult[T any]
// 但这里用最简单的方式，因为结构体不需要泛型
type PageQuery struct {
	Page     int    `form:"page" json:"page"`           // 当前页码，从 1 开始
	PageSize int    `form:"page_size" json:"page_size"` // 每页数量
	Keyword  string `form:"keyword" json:"keyword"`     // 搜索关键词
}

// PageResult 分页结果
type PageResult struct {
	Total int64       `json:"total"` // 总记录数
	List  interface{} `json:"list"`  // 当前页数据
	Page  int         `json:"page"`  // 当前页码
}

// GetPage 从分页参数中获取安全的页码和每页数量
// 【知识点】这种"确保安全值"的函数叫"默认值函数"
// 防止前端传 page=0 或 page=-1 导致异常
func (p *PageQuery) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

func (p *PageQuery) GetPageSize() int {
	if p.PageSize <= 0 {
		return 10
	}
	if p.PageSize > 100 {
		return 100 // 限制最大每页数量，防止一次查太多数据
	}
	return p.PageSize
}

// GetOffset 计算数据库查询的偏移量
// 【知识点】SQL 分页用 LIMIT 和 OFFSET
// 例如第 3 页，每页 10 条: OFFSET = (3-1) * 10 = 20, LIMIT 10
func (p *PageQuery) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}
