package upload

import (
	"time"
)

type Image struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID       uint      `gorm:"column:user_id;not null" json:"user_id"`
	FileName     string    `gorm:"column:file_name;type:varchar(255);not null" json:"file_name"`
	ContentType  string    `gorm:"column:content_type;type:varchar(100);not null" json:"content_type"`
	Size         int64     `gorm:"column:size;not null" json:"size"`
	OriginalName string    `gorm:"column:original_name;type:varchar(255);not null" json:"original_name"`
	UploadIP     string    `gorm:"column:upload_ip;type:varchar(45)" json:"upload_ip"`
	UserAgent    string    `gorm:"column:user_agent;type:varchar(255)" json:"user_agent"`
	UploadedAt   time.Time `gorm:"column:uploaded_at;autoCreateTime" json:"uploaded_at"`
}
