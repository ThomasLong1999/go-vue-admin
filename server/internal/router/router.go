// ============================================
// router.go - 路由注册
// ============================================
// 【知识点】Gin 路由系统: 把 URL 路径映射到处理函数
//
// 路由组（Route Group）: 共享前缀和中间件的一组路由
// 例如: api := r.Group("/api")，注册 "users" 实际匹配 "/api/users"
//
// RESTful API 设计规范:
// GET    /users     → 列表
// POST   /users     → 创建
// GET    /users/:id → 详情
// PUT    /users/:id → 更新
// DELETE /users/:id → 删除
//
// :id 是路径参数，Gin 用 c.Param("id") 获取

package router

import (
	"go-vue-admin/server/internal/handler"
	"go-vue-admin/server/internal/middleware"
	"go-vue-admin/server/internal/pkg"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有路由
func SetupRouter(r *gin.Engine, jwtMgr *pkg.JWTManager, userHandler *handler.UserServiceHandler, dashboardHandler *handler.DashboardHandler) {

	// ===== 公开路由（不需要登录） =====
	api := r.Group("/api")
	{
		api.GET("/health", handler.Health)

		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}
	}

	// ===== 受保护路由（需要 JWT 登录） =====
	// 【知识点】Use() 应用中间件，之后的路由都会经过这个中间件
	// JWT 中间件会验证 Token，通过后把用户信息存入 Context
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth(jwtMgr))
	{
		// 用户管理 CRUD
		users := protected.Group("/users")
		// 登录只证明“你是谁”，管理员检查才决定“你能管理用户”。
		users.Use(middleware.RequireAdmin())
		{
			users.GET("", userHandler.ListUsers)
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// 仪表盘
		dashboard := protected.Group("/dashboard")
		{
			dashboard.GET("/stats", dashboardHandler.DashboardStats)
		}
	}
}
