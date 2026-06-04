// ============================================
// jwt.go - JWT Token 工具包
// ============================================
// 【知识点】JWT（JSON Web Token）是一种无状态的认证方式
// 服务端不需要存储 session，Token 自带用户信息
//
// JWT 由三部分组成，用 . 分隔:
// 1. Header（头部）: 算法类型
// 2. Payload（载荷）: 用户信息 + 过期时间
// 3. Signature（签名）: 用密钥对前两部分签名，防止篡改
//
// 工作流程:
// 用户登录 → 服务端生成 JWT → 前端保存（localStorage）
// → 后续请求在 Header 中携带 JWT → 服务端验证签名和过期时间

package pkg

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 【知识点】常量用 const 声明，编译时确定，不可修改
// 变量用 var 声明，运行时可以改变
const (
	JWTIssuer = "go-vue-admin" // Token 签发者，用于标识是哪个系统发的
)

// Claims 自定义 JWT 载荷
// 【知识点】Go 支持"组合"（embedding），类似其他语言的继承
// jwt.RegisteredClaims 是标准库定义的通用字段
// 把它嵌入 Claims 结构体，Claims 就自动拥有标准字段（过期时间、签发者等）
type Claims struct {
	UserID   uint   `json:"user_id"`   // 用户ID
	Username string `json:"username"` // 用户名
	jwt.RegisteredClaims              // 嵌入标准 JWT 字段
}

// JWTManager JWT 管理器，封装生成和解析逻辑
type JWTManager struct {
	secret     []byte
	expireHour time.Duration
}

// NewJWTManager 创建 JWT 管理器
// 【知识点】构造函数模式：Go 没有构造函数语法，习惯用 NewXxx 函数
func NewJWTManager(secret string, expireHours int) *JWTManager {
	return &JWTManager{
		secret:     []byte(secret),
		expireHour: time.Duration(expireHours) * time.Hour,
	}
}

// GenerateToken 生成 JWT Token
// 参数: userID(用户ID), username(用户名)
// 返回: token 字符串, 错误信息
//
// 【知识点】jwt.NewWithClaims 创建一个新的 Token 对象
// jwt.NewWithClaims(签名算法, 载荷)
// HMAC-SHA256 是最常用的对称加密算法，安全性足够且性能好
func (m *JWTManager) GenerateToken(userID uint, username string) (string, error) {
	// 【知识点】time.Now() 获取当前时间
	// .Add() 方法加上过期时长，得到过期时间
	now := time.Now()

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expireHour)),
			// 【知识点】jwt.NewNumericDate 把 time.Time 转成 JWT 需要的格式
			// JWT 规定时间用 Unix 时间戳（秒）
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    JWTIssuer,
		},
	}

	// 创建 Token 对象并签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// SignedString 用密钥对 Token 签名，生成最终字符串
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("生成 Token 失败: %w", err)
	}

	return tokenString, nil
}

// ParseToken 解析和验证 JWT Token
// 【知识点】jwt.ParseWithClaims 做三件事:
// 1. 解析 Token 字符串
// 2. 验证签名（用密钥重新计算签名，看是否匹配）
// 3. 验证过期时间等标准字段
//
// 如果任何一个验证失败，返回错误
func (m *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	// 【知识点】匿名函数作为参数传入
	// jwt.ParseWithClaims 的第三个参数是一个回调函数
	// 它在验证签名时被调用，需要返回签名用的密钥
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 【知识点】类型断言: token.Method.(jwt.SigningMethodHMAC)
		// Go 是静态类型语言，interface{} 可以是任何类型
		// 使用时需要"断言"把它转成具体类型
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名算法: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析 Token 失败: %w", err)
	}

	// 类型断言获取自定义 Claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的 Token")
}
