package model

import (
	"gorm.io/gorm"
)

type Organization struct {
	Name        string `gorm:"size:100;not null" json:"name"`
	Email       string `gorm:"size:100" json:"email"`
	Avatar      string `gorm:"size:255" json:"avatar"`
	Description string `gorm:"type:text" json:"description"`
	Location    string `gorm:"size:100" json:"location"`
	Website     string `gorm:"size:255" json:"website"`
	gorm.Model
}

type OrganizationMember struct {
	OrganizationID uint   `gorm:"primarykey" json:"organization_id"`
	UserID         uint   `gorm:"primarykey" json:"user_id"`
	Role           string `gorm:"size:20;not null;default:'member'" json:"role"` // admin or member
	gorm.Model
}
