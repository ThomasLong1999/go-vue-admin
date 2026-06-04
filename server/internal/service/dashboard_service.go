// ============================================
// dashboard_service.go - 仪表盘业务逻辑
// ============================================

package service

import (
	"go-vue-admin/server/internal/repository"
)

type DashboardService struct {
	dashboardRepo *repository.DashboardRepository
}

func NewDashboardService(dashboardRepo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{dashboardRepo: dashboardRepo}
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

// GetStats 获取仪表盘统计数据
func (s *DashboardService) GetStats() (*StatsResult, error) {
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

	return result, nil
}
