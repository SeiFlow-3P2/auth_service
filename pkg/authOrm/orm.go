package authOrm

import (
	"github.com/SeiFlow-3P2/auth_service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Kurinov
type AuthOrm struct {
	gorm.DB
}

// CreateUser creates new user in database
func (d *AuthOrm) CreateUser(name string, email string, photoUrl string, telegramId uint, password []byte) error {
	user := domain.User{
		Username:     name,
		Email:        email,
		PhotoUrl:     photoUrl,
		TelegramId:   telegramId,
		PasswordHash: password,
	}
	return d.Create(&user).Error
}

// ChangePassword changes user password
func (d *AuthOrm) ChangePassword(userId uuid.UUID, password []byte) error {
	result := d.Model(&domain.User{ID: userId}).Update("password_hash", password).Update("UpdatedAt", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ChangeEmail changes user username
func (d *AuthOrm) ChangeEmail(userId uuid.UUID, email string) error {
	result := d.Model(&domain.User{ID: userId}).Update("email", email).Update("UpdatedAt", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ChangePhoto changes users photo
func (d *AuthOrm) ChangePhoto(userId uuid.UUID, photoUrl string) error {
	result := d.Model(&domain.User{ID: userId}).Update("photo_url", photoUrl).Update("UpdatedAt", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ChangeTelegramId changes users telegram id
func (d *AuthOrm) ChangeTelegramId(userId uuid.UUID, telegramId uint) error {
	result := d.Model(&domain.User{ID: userId}).Update("telegram_id", telegramId).Update("UpdatedAt", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (d *AuthOrm) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := d.First(&user, email).Error
	return &user, err
}

// GetUser returns user by id
func (d *AuthOrm) GetUser(userId uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := d.First(&user, userId).Error
	return &user, err
}

func (d *AuthOrm) Ping() error {
	db, err := d.DB.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}

func (d *AuthOrm) MigrateDB() error {
	err := d.AutoMigrate(&domain.User{})
	return err
}
