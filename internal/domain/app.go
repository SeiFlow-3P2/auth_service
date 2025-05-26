package domain

import (
	"github.com/SeiFlow-3P2/auth_service/pkg/authRedis"
	authv1 "github.com/SeiFlow-3P2/auth_service/pkg/proto/v1"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type App struct {
	Casher     *authRedis.Casher
	AuthDB     AuthDB
	GrpcServer *authv1.AuthServiceServer
	Settings   *AppSettings
	Logger     *slog.Logger
}
type AppSettings struct {
	Secret     string
	RefreshTTL time.Duration
	AccessTTL  time.Duration
}

type AuthDB interface {
	CreateUser(name string, email string, photoUrl string, telegramId uint, subscription bool, password []byte) error
	ChangePassword(userId uuid.UUID, password []byte) error
	ChangeEmail(userId uuid.UUID, email string) error
	ChangePhoto(userId uuid.UUID, photoUrl string) error
	ChangeTelegramId(userId uuid.UUID, telegramId uint) error
	GetUser(userId uint) (*User, error)
	Ping() error
}
