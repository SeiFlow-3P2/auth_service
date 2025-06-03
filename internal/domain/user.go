package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `gorm:"primaryKey;not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
	Username     string    `gorm:"size:255;uniqueIndex;not null"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	PhotoUrl     string    `gorm:"size:255;default:null"`
	TelegramId   uint      `gorm:"size:11"`
	PasswordHash []byte
}
type UserInfo struct {
	ID        string
	Name      string
	Email     string
	AvatarURL string
	Provider  string
}
