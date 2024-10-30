package model

import (
	"gorm.io/gorm"
	"time"
)

type UserSession struct {
	gorm.Model
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	IsValid   bool      `gorm:"not null;default:true"`
	ExpiredAt time.Time `gorm:"not null"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
