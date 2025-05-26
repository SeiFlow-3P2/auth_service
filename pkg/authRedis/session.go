package authRedis

import (
	"context"
	"errors"
)

// Kurinov

const (
	refreshField = "refresh"
	passField    = "pass"
)

// SetSession sets refresh token and pass hash
func (r *Casher) SetSession(ctx context.Context, userID string, refreshToken string, passHash []byte) error {

	setRef := r.Client.HSet(userID, refreshField, refreshToken)
	if setRef.Err() != nil {
		return setRef.Err()
	}

	setPass := r.Client.HSet(userID, passField, passHash)
	if setPass.Err() != nil {
		return setPass.Err()
	}

	exp := r.Client.Expire(userID, r.RefreshTTL)
	if exp.Err() != nil {
		return exp.Err()

	}
	return nil
}

// UserRefresh returns refresh token
func (r *Casher) UserRefresh(ctx context.Context, userID string) (string, error) {
	var result string

	err := r.Client.HGet(userID, refreshField).Scan(&result)

	if err != nil {
		return "", err
	}
	return result, nil

}

// UserPassHash returns pass hash
func (r *Casher) UserPassHash(ctx context.Context, userID string) (string, error) {
	var result string

	err := r.Client.HGet(userID, passField).Scan(&result)

	if err != nil {
		return "", err
	}
	return result, nil

}

// BlockSession deletes user from casher
func (r *Casher) BlockSession(ctx context.Context, userID string) error {

	err := r.Client.HDel(userID).Err()

	if err != nil {
		return err
	}
	return nil
}

// UserSession returns refresh token and pass hash
func (r *Casher) UserSession(ctx context.Context, userID string) (string, string, error) {
	result, err := r.Client.HGetAll(userID).Result()
	if err != nil {
		return "", "", err
	}

	refresh, ok := result[refreshField]
	if !ok {
		return "", "", errors.New("no refresh token")
	}

	pass, ok := result[passField]
	if !ok {
		return "", "", errors.New("no pass hash")
	}
	return refresh, pass, nil
}
