package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `authOrm:"primaryKey;not null"`
	CreatedAt    time.Time `authOrm:"not null"`
	UpdatedAt    time.Time `authOrm:"not null"`
	Username     string    `authOrm:"size:255;uniqueIndex;not null"`
	Email        string    `authOrm:"type:varchar(100);uniqueIndex;not null"`
	PhotoUrl     string    `authOrm:"size:255;default:null"`
	TelegramId   uint      `authOrm:"size:11"`
	PasswordHash []byte
}
