// ============================================
// handler.go - HTTP 请求处理层（Controller）
// ============================================
// 【知识点】Handler 层负责:
// 1. 从 HTTP 请求中提取参数（路径参数、查询参数、请求体）
// 2. 调用 Service 层处理业务
// 3. 将结果包装成统一响应格式返回
//
// Handler 不应该包含业务逻辑，它只是"翻译官"：
// 把 HTTP 请求翻译成 Service 调用，把 Service 结果翻译成 HTTP 响应

package handler

import (
	"strconv"

	"go-vue-admin/server/internal/model"
	"go-vue-admin/server/internal/pkg"
	"go-vue-admin/server/internal/service"

	"github.com/gin-gonic/gin"
)

// UserServiceHandler 用户相关 HTTP 处理函数
type UserServiceHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserServiceHandler {
	return &UserServiceHandler{userSvc: userSvc}
}

// Health 健康检查接口
// 【知识点】Handler 函数签名固定: func(*gin.Context)
// gin.Context 是 Gin 的核心对象，贯穿整个请求生命周期
func Health(c *gin.Context) {
	pkg.Success(c, gin.H{
		"status":  "running",
		"message": "go-vue-admin 服务正常运行",
	})
	// 【知识点】gin.H 是 map[string]interface{} 的快捷方式
	// 用于快速构建 JSON 响应
}

// Register 用户注册
func (h *UserServiceHandler) Register(c *gin.Context) {
	var input service.RegisterInput
	// 【知识点】ShouldBindJSON 从请求体解析 JSON 到结构体
	// 同时校验 binding tag（required, min, max, email 等）
	if err := c.ShouldBindJSON(&input); err != nil {
		pkg.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	user, err := h.userSvc.Register(input)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, user)
}

// Login 用户登录
func (h *UserServiceHandler) Login(c *gin.Context) {
	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		pkg.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	user, token, err := h.userSvc.Login(input)
	if err != nil {
		pkg.Fail(c, 401, err.Error())
		return
	}

	pkg.Success(c, gin.H{
		"user":  user,
		"token": token,
	})
}

// ListUsers 获取用户列表
// 【知识点】ShouldBindQuery 从 URL 查询参数解析
// GET /api/users?page=1&page_size=10&keyword=abc
func (h *UserServiceHandler) ListUsers(c *gin.Context) {
	var query model.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		pkg.Fail(c, 400, "参数错误")
		return
	}

	users, total, err := h.userSvc.ListUsers(query)
	if err != nil {
		pkg.Fail(c, 500, "查询失败")
		return
	}

	pkg.Success(c, model.PageResult{
		Total: total,
		List:  users,
		Page:  query.GetPage(),
	})
}

// GetUser 获取单个用户
// 【知识点】c.Param("id") 获取 URL 路径参数
// 路由: GET /users/:id → 访问: GET /users/123 → Param("id") = "123"
func (h *UserServiceHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Fail(c, 400, "无效的用户ID")
		return
	}

	user, err := h.userSvc.GetUser(uint(id))
	if err != nil {
		pkg.Fail(c, 404, "用户不存在")
		return
	}

	pkg.Success(c, user)
}

// CreateUser 管理员创建用户
func (h *UserServiceHandler) CreateUser(c *gin.Context) {
	var input service.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		pkg.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	user, err := h.userSvc.Register(input)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, user)
}

// UpdateUser 更新用户
func (h *UserServiceHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Fail(c, 400, "无效的用户ID")
		return
	}

	var updates map[string]interface{}
	// 【知识点】绑定到 map 可以接受任意字段，适合更新操作
	if err := c.ShouldBindJSON(&updates); err != nil {
		pkg.Fail(c, 400, "参数错误")
		return
	}

	user, err := h.userSvc.UpdateUser(uint(id), updates)
	if err != nil {
		pkg.Fail(c, 500, "更新失败")
		return
	}

	pkg.Success(c, user)
}

// DeleteUser 删除用户
func (h *UserServiceHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Fail(c, 400, "无效的用户ID")
		return
	}

	if err := h.userSvc.DeleteUser(uint(id)); err != nil {
		pkg.Fail(c, 500, "删除失败")
		return
	}

	pkg.Success(c, nil)
}

// ============================================
// DashboardHandler 仪表盘处理器
// ============================================
type DashboardHandler struct {
	dashboardSvc *service.DashboardService
}

func NewDashboardHandler(dashboardSvc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardSvc: dashboardSvc}
}

func (h *DashboardHandler) DashboardStats(c *gin.Context) {
	stats, err := h.dashboardSvc.GetStats()
	if err != nil {
		pkg.Fail(c, 500, "获取统计数据失败")
		return
	}
	pkg.Success(c, stats)
}
