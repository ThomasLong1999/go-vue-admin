// ============================================
// jwt.go - JWT 鉴权中间件
// ============================================
// 【知识点】中间件（Middleware）是 Gin 框架的核心特性
// 它是一个函数，在请求到达 Handler 之前/之后执行
// 类比: 快递分拣中心的安检环节，每个包裹都要过安检
//
// Gin 中间件的执行流程（洋葱模型）:
// 请求 → 中间件前置代码 → Handler → 中间件后置代码 → 响应
//
// 使用方式:
// router.Use(middleware.JWTAuth(jwtMgr))
// 之后所有注册的路由都会经过这个中间件

package middleware

import (
	"strings"

	"go-vue-admin/server/internal/pkg"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 鉴权中间件
// 【知识点】闭包（Closure）: 外层函数 JWTAuth(jwtMgr) 返回一个内层函数
// 内层函数可以"捕获"并使用外层函数的参数 jwtMgr
func JWTAuth(jwtMgr *pkg.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取 Authorization
		// 【知识点】标准格式: Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			pkg.FailWithHTTPCode(c, 401, 401, "未登录，请先登录")
			c.Abort() // 终止请求链，不调 Abort 的话后续 Handler 仍会执行
			return
		}

		// 2. 提取 Token（去掉 "Bearer " 前缀）
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			pkg.FailWithHTTPCode(c, 401, 401, "Token 格式错误")
			c.Abort()
			return
		}

		// 3. 解析和验证 Token
		claims, err := jwtMgr.ParseToken(tokenString)
		if err != nil {
			pkg.FailWithHTTPCode(c, 401, 401, "Token 无效或已过期")
			c.Abort()
			return
		}

		// 4. 将用户信息存入 Context，后续 Handler 可通过 c.Get("user_id") 取出
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.IsAdmin)

		c.Next()
	}
}

// RequireAdmin 演示“已登录”和“有权限”是两个不同的概念。
// JWTAuth 负责确认身份；本中间件只允许管理员继续访问。
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !c.GetBool("is_admin") {
			pkg.FailWithHTTPCode(c, 403, 403, "需要管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}

// CORS 跨域中间件
// 【知识点】CORS（Cross-Origin Resource Sharing）跨域资源共享
// 浏览器同源策略: 协议+域名+端口 三个都相同才算同源
// 前端（localhost:5173）调后端（localhost:8080）是跨域请求
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		// OPTIONS 是浏览器的"预检请求"
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
