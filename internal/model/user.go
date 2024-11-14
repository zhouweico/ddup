package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint       `gorm:"primarykey"`
	Username      string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password      string     `gorm:"type:varchar(100);not null" json:"-"`
	Nickname      string     `gorm:"type:varchar(50);not null" json:"nickname"`
	Gender        string     `gorm:"size:10;default:'unknown'" json:"gender"`
	Birthday      *time.Time `json:"birthday"`
	Avatar        string     `gorm:"type:varchar(255)" json:"avatar"`
	Email         string     `gorm:"size:100;null" json:"email"`
	Mobile        string     `gorm:"size:20;null" json:"mobile"`
	Location      string     `gorm:"size:100;null" json:"location"`
	Language      string     `gorm:"type:varchar(10);default:zh-CN" json:"language"`
	Bio           string     `gorm:"size:500" json:"bio"`
	Status        int        `gorm:"default:1;not null" json:"status"`
	LastLogin     *time.Time `json:"last_login"`
	LoginAttempts int        `gorm:"default:0" json:"-"`
	LockedUntil   *time.Time `json:"-"`
	gorm.Model
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Nickname == "" {
		u.Nickname = u.Username
	}
	if u.Gender == "" {
		u.Gender = "unknown"
	}
	return nil
}
