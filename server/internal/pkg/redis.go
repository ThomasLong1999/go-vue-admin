// ============================================
// redis.go - Redis 连接初始化
// ============================================
// 【知识点】go-redis 是 Go 生态最流行的 Redis 客户端
// 它支持 context（上下文），可以设置超时和取消

package pkg

import (
	"context"
	"fmt"
	"log"

	"go-vue-admin/server/config"

	"github.com/redis/go-redis/v9"
)

// RDB 是全局 Redis 连接对象
var RDB *redis.Client

// InitRedis 初始化 Redis 连接
// 【知识点】context.Background() 创建一个空上下文
// context 是 Go 并发编程的核心概念，用于:
// 1. 传递请求级别的值（如 traceID）
// 2. 控制超时和取消
// Background() 表示"没有截止时间，也不可取消"的根上下文
func InitRedis(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 【知识点】Ping 检查 Redis 是否连接成功
	// 这一步很重要：创建客户端不代表连接成功
	// Ping 发送一个命令测试连接，失败就说明配置有误或 Redis 没启动
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("连接 Redis 失败: %w", err)
	}

	log.Printf("Redis 连接成功: %s", cfg.Addr())

	RDB = rdb
	return rdb, nil
}
