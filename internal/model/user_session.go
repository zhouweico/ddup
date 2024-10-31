package model

import (
	"time"

	"gorm.io/gorm"
)

type UserSession struct {
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	IsValid   bool      `gorm:"not null;default:true"`
	ExpiredAt time.Time `gorm:"not null"`
	gorm.Model
}

func (UserSession) TableName() string {
	return "user_sessions"
}
