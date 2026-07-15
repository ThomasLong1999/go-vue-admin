// ============================================
// dashboard_service.go - 仪表盘业务逻辑
// ============================================

package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"go-vue-admin/server/internal/repository"

	"github.com/redis/go-redis/v9"
)

type DashboardService struct {
	dashboardRepo *repository.DashboardRepository
	cache         *redis.Client
}

func NewDashboardService(dashboardRepo *repository.DashboardRepository, cache *redis.Client) *DashboardService {
	return &DashboardService{dashboardRepo: dashboardRepo, cache: cache}
}

// StatsResult 仪表盘统计数据
type StatsResult struct {
	TotalUsers    int64        `json:"total_users"`
	TodayNew      int64        `json:"today_new"`
	ActiveUsers   int64        `json:"active_users"`
	DisabledUsers int64        `json:"disabled_users"`
	RecentUsers   []RecentUser `json:"recent_users"`
}

type RecentUser struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	CreatedAt string `json:"created_at"`
}

// GetStats 获取仪表盘统计数据。
// 这里使用最容易理解的 cache-aside 模式：先读缓存，未命中再读数据库。
func (s *DashboardService) GetStats(ctx context.Context) (*StatsResult, error) {
	const cacheKey = "dashboard:stats"
	if s.cache != nil {
		cached, err := s.cache.Get(ctx, cacheKey).Bytes()
		if err == nil {
			var result StatsResult
			if err := json.Unmarshal(cached, &result); err == nil {
				return &result, nil
			} else {
				log.Printf("仪表盘缓存解析失败，将查询数据库: %v", err)
			}
		} else if !errors.Is(err, redis.Nil) {
			log.Printf("读取仪表盘缓存失败，将查询数据库: %v", err)
		}
	}

	stats, err := s.dashboardRepo.GetUserStats()
	if err != nil {
		return nil, err
	}

	users, err := s.dashboardRepo.RecentUsers(5)
	if err != nil {
		return nil, err
	}

	result := &StatsResult{
		TotalUsers:    stats.TotalUsers,
		TodayNew:      stats.TodayNew,
		ActiveUsers:   stats.ActiveUsers,
		DisabledUsers: stats.DisabledUsers,
	}

	// 转换最近用户数据
	for _, u := range users {
		result.RecentUsers = append(result.RecentUsers, RecentUser{
			ID:        u.ID,
			Username:  u.Username,
			Nickname:  u.Nickname,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04"),
			// 【知识点】Go 的时间格式化很特别
			// 不是用 YYYY-MM-DD，而是用 2006-01-02 15:04:05
			// 原因是 Go 诞生于 2006年1月2日 15:04:05，用这个做参考
		})
	}

	if s.cache != nil {
		payload, err := json.Marshal(result)
		if err != nil {
			log.Printf("仪表盘缓存序列化失败: %v", err)
		} else if err := s.cache.Set(ctx, cacheKey, payload, 30*time.Second).Err(); err != nil {
			// 缓存只是加速层，失败时仍把数据库结果返回给页面。
			log.Printf("写入仪表盘缓存失败: %v", err)
		}
	}

	return result, nil
}
