package service

import (
	"context"
	"errors"
	"github.com/SeiFlow-3P2/auth_service/internal/domain"
	"github.com/SeiFlow-3P2/auth_service/pkg/authJWT"
	verfic "github.com/SeiFlow-3P2/auth_service/pkg/utils/verifications"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type Auth struct {
	*domain.App
}

func (a *Auth) SingUpByOauth(ctx context.Context, provider string, oauthToken string, telegramID string) (userID uuid.UUID, accessToken string, refreshToken string, message string, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *Auth) LoginByEmail(ctx context.Context, email string, password []byte) (userID uuid.UUID, accessToken string, refreshToken string, message string, err error) {

	refreshToken, passwordHash, err := a.Casher.UserSession(ctx, email)
	if refreshToken == "" {
		user, err := a.AuthDB.GetUserByEmail(email)
		if err != nil {
			return uuid.Nil, "", "", "", err
		}
		if user == nil {
			return uuid.Nil, "", "", "", errors.New("user not found")
		}
		if string(user.PasswordHash) != string(password) {
			return uuid.Nil, "", "", "", errors.New("invalid password")
		}
		tokens, err := authJWT.CreateTokenPair(ctx, *user, a.Settings)
		if err != nil {
			return uuid.Nil, "", "", "", err
		}

		return user.ID, tokens.AccessToken, tokens.RefreshToken, "", err
	}

	if err != nil {
		return uuid.Nil, "", "", "", err
	}

	if string(passwordHash) != string(password) {
		return uuid.Nil, "", "", "", errors.New("invalid password")
	}

	accessToken, refreshToken, err = a.RefreshToken(ctx, refreshToken)

	if err != nil {
		return uuid.Nil, "", "", "", err
	}
	return userID, accessToken, refreshToken, "", err

}

func (a *Auth) LoginByOauth(ctx context.Context, provider string, oauthToken string) (userID uuid.UUID, accessToken string, refreshToken string, message string, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *Auth) RefreshToken(ctx context.Context, RefreshToken string) (accessToken string, refreshToken string, err error) {

	refToken, err := jwt.Parse(RefreshToken,
		func(refToken *jwt.Token) (interface{}, error) {
			tokenChecked, err := refToken.SignedString(a.Settings.Secret)
			if err != nil {
				return nil, err
			}
			return tokenChecked, nil
		})

	if err != nil {
		return "", "", err

	}
	claims, ok := refToken.Claims.(jwt.MapClaims)
	if ok != true {
		return "", "", errors.New("invalid token")
	}

	expired := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(expired) {
		return "", "", errors.New("token expired")
	}

	userID, err := uuid.Parse(claims["id"].(string))

	user, err := a.AuthDB.GetUser(userID)
	if err != nil || user == nil {
		return "", "", err
	}

	tokens, err := authJWT.CreateTokenPair(ctx, *user, a.Settings)
	if err != nil {
		return "", "", err
	}
	err = a.Casher.BlockSession(ctx, user.Email)

	err = a.Casher.SetSession(ctx, user.Email, tokens.RefreshToken, user.PasswordHash)
	return tokens.AccessToken, tokens.RefreshToken, nil

}

func (a *Auth) Logout(ctx context.Context, userID uuid.UUID) (err error) {
	err = a.Casher.BlockSession(ctx, userID.String())
	return err
}

func (a *Auth) SingUpByEmail(ctx context.Context, name string, email string, password []byte, telegramID uint) (userID uuid.UUID, accessToken string, refreshToken string, message string, err error) {
	validateEmail, err := verfic.VerifyEmail(email)
	if err != nil || !validateEmail {
		return uuid.Nil, "", "", "Неверный формат email", err
	}
	err = a.AuthDB.CreateUser(name, email, "", telegramID, password)
	if err != nil {
		return uuid.Nil, "", "", "", err
	}
	userID, accessToken, refreshToken, message, err = a.LoginByEmail(ctx, email, password)
	if err != nil {
		return uuid.Nil, "", "", "", err
	}
	return userID, accessToken, refreshToken, message, nil
}

func (a *Auth) UserInfo(ctx context.Context, userID uuid.UUID) (id string, telegramID uint, username string, email string, photoUrl string, createdAt string, updatedAt string, err error) {
	user, err := a.AuthDB.GetUser(userID)
	if err != nil {
		return "", 0, "", "", "", "", "", err
	}
	if user != nil {
		return user.ID.String(), user.TelegramId, user.Username, user.Email, user.PhotoUrl, user.CreatedAt.String(), user.UpdatedAt.String(), nil
	}
	return "", 0, "", "", "", "", "", errors.New("не удалось получить пользователя")
}

func (a *Auth) HealthCheck(ctx context.Context) (status string, err error) {

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
