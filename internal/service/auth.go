package service

import (
	"context"
	"github.com/SeiFlow-3P2/auth_service/internal/domain"
	"github.com/SeiFlow-3P2/auth_service/pkg/authRedis"
	"github.com/google/uuid"
	"log/slog"
)

type Auth struct {
	AuthDB domain.AuthDB
	Casher *authRedis.Casher
	Logger *slog.Logger
}

func (a Auth) SingUpByEmail(ctx context.Context, email string, password string, telegramID string) (userID uuid.UUID, accessToken string, refreshToken string, message string, err error) {
	//TODO implement me
	panic("implement me")
}

func (a Auth) SingUpByOauth(ctx context.Context, provider string, oauthToken string, telegramID string) (userID uuid.UUID, accessToken string, refreshToken string, message string, err error) {
	//TODO implement me
	panic("implement me")
}

func (a Auth) LoginByEmail(ctx context.Context, email string, password string, telegramID string) (userID uuid.UUID, accessToken string, refreshToken string, message string, err error) {
	//TODO implement me
	panic("implement me")
}

func (a Auth) LoginByOauth(ctx context.Context, provider string, oauthToken string) (userID uuid.UUID, accessToken string, refreshToken string, message string, err error) {
	//TODO implement me
	panic("implement me")
}

func (a Auth) RefreshToken(ctx context.Context, RefreshToken string) (accessToken string, refreshToken string, err error) {
	//TODO implement me
	panic("implement me")
}

func (a Auth) Logout(ctx context.Context, RefreshToken string) {
	//TODO implement me
	panic("implement me")
}

func (a Auth) GetUserInfo(ctx context.Context, userID uuid.UUID) (id string, telegramId string, username string, email string, photoUrl string, subscription bool, createdAt string, updatedAt string, err error) {
	//TODO implement me
	panic("implement me")
}

func (a Auth) HealthCheck(ctx context.Context) (status string, err error) {

	err = a.AuthDB.Ping()
	if err != nil {
		a.Logger.Error("Ошибка подключения к базе данных", slog.Any("err", err))
		return "FAIL", err
	}

	err = a.Casher.Ping().Err()
	if err != nil {
		a.Logger.Error("Ошибка подключения к кешу", slog.Any("err", err))
		return "FAIL", err
	}
	return "OK", nil
}
