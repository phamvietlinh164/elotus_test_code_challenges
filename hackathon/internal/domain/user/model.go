package user

import "time"

type User struct {
	ID             uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username       string    `gorm:"column:username;type:varchar(100);uniqueIndex;not null" json:"username"`
	Name           string    `gorm:"column:name;type:varchar(255)" json:"name"`
	Email          string    `gorm:"column:email;type:varchar(255);uniqueIndex" json:"email"`
	Password       string    `gorm:"column:password;type:varchar(255);not null" json:"-"`
	IsAdmin        bool      `gorm:"column:is_admin;default:false" json:"is_admin"`
	TokenRevokedAt time.Time `gorm:"column:token_revoked_at;default:null" json:"token_revoked_at"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}
