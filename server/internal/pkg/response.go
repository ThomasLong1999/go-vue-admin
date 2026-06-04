// ============================================
// response.go - 统一响应格式封装
// ============================================
// 【知识点】在真实项目中，所有 API 接口的返回格式应该统一
// 前端只需要按一种格式解析响应，不需要每个接口单独处理
//
// 标准格式:
// {
//   "code": 0,        // 业务状态码，0 表示成功
//   "message": "ok",  // 提示信息
//   "data": {...}     // 实际数据
// }

package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 业务状态码: 0=成功, 非0=错误
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 实际数据
	// 【知识点】interface{} 是 Go 的"空接口"，可以表示任何类型
	// 类似 Java 的 Object、TypeScript 的 unknown
	// 用它是因为 data 字段可能是用户对象、列表、或 null
}

// Success 返回成功响应
// 【知识点】gin.Context 是 Gin 框架中最重要的对象
// 它贯穿整个请求生命周期，包含请求信息、响应方法、中间件数据等
//
// 【知识点】*gin.Context 中的 * 表示指针
// 指针传递的是对象的地址，函数内修改会影响原对象
// 如果传值（不加 *），修改的是副本，不影响原对象
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

// SuccessWithMessage 返回带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Fail 返回失败响应
func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// FailWithHTTPCode 返回带 HTTP 错误码的失败响应
// 【知识点】HTTP 状态码（如 404, 500）和业务状态码是两套体系
// HTTP 状态码是协议层面的，客户端/浏览器会根据它做不同处理
// 业务状态码是应用层面的，前端根据它来判断业务是否成功
func FailWithHTTPCode(c *gin.Context, httpCode int, code int, message string) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
