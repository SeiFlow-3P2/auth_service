package authJWT

import (
	"context"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"

	"github.com/SeiFlow-3P2/auth_service/internal/domain"
)

func createAccessToken(ctx context.Context, User domain.User, tokenUUID uuid.UUID, Settings *domain.AppSettings) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["uuid"] = tokenUUID
	claims["username"] = User.Username
	claims["email"] = User.Email
	claims["telegram_id"] = User.TelegramId
	claims["created_at"] = User.CreatedAt
	claims["updated_at"] = User.UpdatedAt
	claims["exp"] = time.Now().Add(Settings.AccessTTL).Unix()

	secret, err := base64.StdEncoding.DecodeString(Settings.Secret)
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func createRefreshToken(ctx context.Context, User domain.User, tokenUUID uuid.UUID, Settings *domain.AppSettings) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["uuid"] = tokenUUID
	claims["email"] = User.Email
	claims["exp"] = time.Now().Add(Settings.RefreshTTL).Unix()

	secret, err := base64.StdEncoding.DecodeString(Settings.Secret)
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateTokenPair(ctx context.Context, User domain.User, Settings *domain.AppSettings) (domain.Tokens, error) {
	tokensUUID, err := uuid.NewRandom()
	if err != nil {
		return domain.Tokens{}, err
	}

	accessToken, err := createAccessToken(ctx, User, tokensUUID, Settings)
	if err != nil {
		return domain.Tokens{}, err
	}

	refreshToken, err := createRefreshToken(ctx, User, tokensUUID, Settings)
	if err != nil {
		return domain.Tokens{}, err
	}

	return domain.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
