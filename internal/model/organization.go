package model

import (
	"gorm.io/gorm"
)

type Organization struct {
	ID          uint   `gorm:"primarykey"`
	Name        string `gorm:"size:100;uniqueIndex;not null;check:name ~* '^[a-z0-9][a-z0-9-]{0,38}[a-z0-9]$'" json:"name"`
	DisplayName string `gorm:"size:100;not null" json:"display_name"`
	Email       string `gorm:"size:100" json:"email"`
	Avatar      string `gorm:"size:255" json:"avatar"`
	Description string `gorm:"type:text" json:"description"`
	Location    string `gorm:"size:100" json:"location"`
	Website     string `gorm:"size:255" json:"website"`
	gorm.Model
}

type OrganizationMember struct {
	ID             uint   `gorm:"primarykey"`
	OrganizationID uint   `gorm:"primarykey" json:"organization_id"`
	UserID         uint   `gorm:"primarykey" json:"user_id"`
	Role           string `gorm:"size:20;not null;default:'member'" json:"role"` // admin or member
	gorm.Model
}
