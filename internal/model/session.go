package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"type:varchar(500);not null"`
	IsValid   bool      `gorm:"not null;default:true"`
	ExpiredAt time.Time `gorm:"not null"`
	gorm.Model
}
