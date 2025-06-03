package service

import (
	"context"
	"encoding/base64"
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

		err = a.Casher.SetSession(ctx, user.Email, tokens.RefreshToken, user.PasswordHash)

		return user.ID, tokens.AccessToken, tokens.RefreshToken, "", err
	}

	if err != nil {
		return uuid.Nil, "", "", "", err
	}

	if string(passwordHash) != string(password) {
		return uuid.Nil, "", "", "", errors.New("invalid password")
	}

	user, err := a.AuthDB.GetUserByEmail(email)

	if err != nil {

		return uuid.Nil, "", "", "", err
	}

	tokens, err := authJWT.CreateTokenPair(ctx, *user, a.Settings)

	if err != nil {
		return uuid.Nil, "", "", "", err
	}
	return user.ID, tokens.AccessToken, tokens.RefreshToken, "", err

}

func (a *Auth) LoginByOauth(ctx context.Context, provider string, oauthToken string) (userID uuid.UUID, accessToken string, refreshToken string, message string, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *Auth) RefreshToken(ctx context.Context, RefreshToken string) (accessToken string, refreshToken string, err error) {
	op := "Auth_Service_RefreshToken: "
	refToken, err := jwt.Parse(RefreshToken,
		func(token *jwt.Token) (interface{}, error) {
			secret, err := base64.StdEncoding.DecodeString(a.Settings.Secret)
			if err != nil {
				a.Logger.Info(op + err.Error())
				return nil, err
			}
			return secret, nil
		})

	if err != nil {
		a.Logger.Error(op, slog.String(op, err.Error()))
		return "", "", err

	}
	claims, ok := refToken.Claims.(jwt.MapClaims)
	if ok != true {
		a.Logger.Error(op, slog.String(op, "cant get claims"))
		return "", "", errors.New("invalid token")
	}

	expired := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(expired) {
		return "", "", errors.New("token expired")
	}

	email := claims["email"].(string)

	refreshToken, _, err = a.Casher.UserSession(ctx, email)

	if err != nil {
		return "", "", err
	}
	if refreshToken == "" {
		return "", "", errors.New("invalid token")
	}
	if refreshToken != RefreshToken {
		return "", "", errors.New("invalid token")
	}

	user, err := a.AuthDB.GetUserByEmail(email)

	if err != nil {
		return "", "", err
	}
	tokens, err := authJWT.CreateTokenPair(ctx, *user, a.Settings)
	if err != nil {
		a.Logger.Info(err.Error())
		return "", "", err
	}
	err = a.Casher.BlockSession(ctx, user.Email)
	err = a.Casher.SetSession(ctx, user.Email, tokens.RefreshToken, user.PasswordHash)

	return tokens.AccessToken, tokens.RefreshToken, nil

}

func (a *Auth) Logout(ctx context.Context, userID uuid.UUID) (err error) {
	user, err := a.AuthDB.GetUser(userID)

	if err != nil {
		return err
	}
	err = a.Casher.BlockSession(ctx, user.Email)

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
