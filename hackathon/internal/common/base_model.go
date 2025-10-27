package common

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string         `gorm:"type:char(36);primaryKey" json:"ID"` // Store UUID as CHAR(36)
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate is a GORM hook that runs before a new record is created
func (baseModel *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if baseModel.ID == "" {
		baseModel.ID = uuid.New().String() // Generate UUID as a string
	}
	baseModel.CreatedAt = time.Now().UTC()
	baseModel.UpdatedAt = time.Now().UTC()
	return nil
}

func (baseModel *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	baseModel.UpdatedAt = time.Now().UTC()
	return
}
