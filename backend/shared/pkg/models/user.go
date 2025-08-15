package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint64         `gorm:"primarykey" json:"id"`
	Email     string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Nickname  string         `gorm:"size:50;not null" json:"nickname"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Status    UserStatus     `gorm:"type:tinyint;default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserStatus 用户状态
type UserStatus int

const (
	UserStatusInactive UserStatus = 0 // 未激活
	UserStatusActive   UserStatus = 1 // 活跃
	UserStatusBanned   UserStatus = 2 // 被禁用
)

// TableName 表名
func (User) TableName() string {
	return "users"
}
