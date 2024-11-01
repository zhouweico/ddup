package model

import "gorm.io/gorm"

// Social 社交媒体账号
type Social struct {
	UserID      uint   `gorm:"not null;index"`             // 用户ID
	Platform    string `gorm:"type:varchar(50);not null"`  // 平台名称
	Username    string `gorm:"type:varchar(100);not null"` // 平台用户名
	URL         string `gorm:"type:varchar(255)"`          // 个人主页链接
	Verified    bool   `gorm:"default:false"`              // 是否已验证
	Description string `gorm:"type:text"`                  // 描述
	gorm.Model

	User User `gorm:"foreignkey:UserID"` // 关联用户
}
