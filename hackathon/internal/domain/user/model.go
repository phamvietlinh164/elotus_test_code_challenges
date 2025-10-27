package user

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Username       string    `gorm:"uniqueIndex;size:100" json:"username"`
	Name           string    `json:"name"`
	Email          string    `gorm:"unique" json:"email"`
	Password       string    `json:"-"`
	IsAdmin        bool      `gorm:"default:false" json:"is_admin"`
	TokenRevokedAt time.Time `gorm:"column:token_revoked_at" json:"token_revoked_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
