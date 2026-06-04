// ============================================
// dashboard_repository.go - 仪表盘数据查询
// ============================================

package repository

import (
	"time"

	"go-vue-admin/server/internal/model"

	"gorm.io/gorm"
)

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

// UserStats 用户统计数据
type UserStats struct {
	TotalUsers   int64 `json:"total_users"`
	TodayNew     int64 `json:"today_new"`
	ActiveUsers  int64 `json:"active_users"`
	DisabledUsers int64 `json:"disabled_users"`
}

// GetUserStats 获取用户统计
func (r *DashboardRepository) GetUserStats() (*UserStats, error) {
	var stats UserStats

	// 总用户数
	r.db.Model(&model.User{}).Count(&stats.TotalUsers)

	// 今日新增
	today := time.Now().Truncate(24 * time.Hour)
	// 【知识点】Truncate 截断到指定精度，这里截断到天，得到今天零点
	r.db.Model(&model.User{}).Where("created_at >= ?", today).Count(&stats.TodayNew)

	// 启用/禁用用户数
	r.db.Model(&model.User{}).Where("status = ?", 1).Count(&stats.ActiveUsers)
	r.db.Model(&model.User{}).Where("status = ?", 0).Count(&stats.DisabledUsers)

	return &stats, nil
}

// RecentUsers 获取最近注册的用户
func (r *DashboardRepository) RecentUsers(limit int) ([]model.User, error) {
	var users []model.User
	err := r.db.Order("created_at DESC").Limit(limit).Find(&users).Error
	return users, err
}
