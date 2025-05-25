package gorm

import (
	"gorm.io/gorm"
	"time"
)

// Kurinov
type User struct {
	ID           uint `gorm:"primaryKey;autoIncrement;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Username     string `gorm:"size:255;uniqueIndex;not null"`
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null"`
	PhotoUrl     string `gorm:"size:255;default:null"`
	TelegramId   int64  `gorm:"size:11"`
	Subscription bool   `gorm:"default:false"`
	PasswordHash []byte
}

// CreateUser creates new user in database
func CreateUser(db *gorm.DB, name string, email string, photoUrl string, telegramId int64, subscription bool, password []byte) error {
	user := User{
		Username:     name,
		Email:        email,
		PhotoUrl:     photoUrl,
		TelegramId:   telegramId,
		Subscription: subscription,
		PasswordHash: password,
	}
	return db.Create(&user).Error
}

// UpdateSubscription updates user subscription
func UpdateSubscription(db *gorm.DB, userId uint) error {
	result := db.Model(&User{ID: userId}).Update("subscription", true).Update("UpdatedAt", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ChangePassword changes user password
func ChangePassword(db *gorm.DB, userId uint, password []byte) error {
	result := db.Model(&User{ID: userId}).Update("password_hash", password).Update("UpdatedAt", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ChangeEmail changes user username
func ChangeEmail(db *gorm.DB, userId uint, email string) error {
	result := db.Model(&User{ID: userId}).Update("email", email).Update("UpdatedAt", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ChangePhoto changes users photo
func ChangePhoto(db *gorm.DB, userId uint, photoUrl string) error {
	result := db.Model(&User{ID: userId}).Update("photo_url", photoUrl).Update("UpdatedAt", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ChangeTelegramId changes users telegram id
func ChangeTelegramId(db *gorm.DB, userId uint, telegramId int64) error {
	result := db.Model(&User{ID: userId}).Update("telegram_id", telegramId).Update("UpdatedAt", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetUser returns user by id
func GetUser(db *gorm.DB, userId uint) (*User, error) {
	var user User
	err := db.First(&user, userId).Error
	return &user, err
}
