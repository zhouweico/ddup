package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// ProfileType 个人资料类型
type ProfileType string

const (
	General       ProfileType = "general"
	Project       ProfileType = "project"
	SideProject   ProfileType = "side_project"
	Exhibition    ProfileType = "exhibition"
	Speaking      ProfileType = "speaking"
	Writing       ProfileType = "writing"
	Award         ProfileType = "award"
	Feature       ProfileType = "feature"
	Work          ProfileType = "work"
	Volunteering  ProfileType = "volunteering"
	Education     ProfileType = "education"
	Certification ProfileType = "certification"
	Contact       ProfileType = "contact"
	Team          ProfileType = "team"
)

// Profile 基础模型
type Profile struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	UserID       uint            `json:"user_id" gorm:"not null;index"`
	Type         ProfileType     `json:"type" gorm:"type:varchar(20);not null"`
	Title        string          `json:"title" gorm:"type:varchar(100);not null"`
	Year         *int            `json:"year" gorm:""`
	StartDate    *time.Time      `json:"start_date" gorm:""`
	EndDate      *time.Time      `json:"end_date" gorm:""`
	Organization string          `json:"organization" gorm:"type:varchar(100)"`
	Location     string          `json:"location" gorm:"type:varchar(100)"`
	URL          string          `json:"url" gorm:"type:varchar(255)"`
	Description  string          `json:"description" gorm:"type:text"`
	Metadata     json.RawMessage `json:"metadata" gorm:"type:jsonb"`
	DisplayOrder int             `json:"display_order" gorm:"default:0"`
	Visibility   string          `json:"visibility" gorm:"type:varchar(10);default:public;check:visibility in ('public','private')"`
	gorm.Model
}

// ProfileMetadata 元数据结构
type ProfileMetadata struct {
	// General
	DisplayName string `json:"display_name,omitempty"`
	WhatYouDo   string `json:"what_you_do,omitempty"`
	Pronouns    string `json:"pronouns,omitempty"`
	About       string `json:"about,omitempty"`

	// Project/SideProject
	Client        string   `json:"client,omitempty"`
	Collaborators []string `json:"collaborators,omitempty"`

	// Work/Education
	Degree    string   `json:"degree,omitempty"`
	Title     string   `json:"title,omitempty"`
	Coworkers []string `json:"coworkers,omitempty"`

	// Contact
	Platform     string `json:"platform,omitempty"`
	Username     string `json:"username,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	CustomName   string `json:"custom_name,omitempty"`

	// Certification
	IssueDate  *time.Time `json:"issue_date,omitempty"`
	ExpiryDate *time.Time `json:"expiry_date,omitempty"`

	// Common
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment 附件结构
type Attachment struct {
	Type     string `json:"type"`      // page/media
	URL      string `json:"url"`       // 文件链接
	MimeType string `json:"mime_type"` // 文件类型
	Size     int64  `json:"size"`      // 文件大小
	Name     string `json:"name"`      // 文件名
}
