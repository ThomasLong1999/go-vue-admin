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
	TotalUsers    int64 `json:"total_users"`
	TodayNew      int64 `json:"today_new"`
	ActiveUsers   int64 `json:"active_users"`
	DisabledUsers int64 `json:"disabled_users"`
}

// GetUserStats 获取用户统计
func (r *DashboardRepository) GetUserStats() (*UserStats, error) {
	var stats UserStats

	// 总用户数
	if err := r.db.Model(&model.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// 今日新增
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// 使用本地日期构造零点，避免 Truncate(24h) 在东八区得到早上 8 点。
	if err := r.db.Model(&model.User{}).Where("created_at >= ?", today).Count(&stats.TodayNew).Error; err != nil {
		return nil, err
	}

	// 启用/禁用用户数
	if err := r.db.Model(&model.User{}).Where("status = ?", 1).Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.User{}).Where("status = ?", 0).Count(&stats.DisabledUsers).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// RecentUsers 获取最近注册的用户
func (r *DashboardRepository) RecentUsers(limit int) ([]model.User, error) {
	var users []model.User
	err := r.db.Order("created_at DESC").Limit(limit).Find(&users).Error
	return users, err
}
