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
func (r *Casher) SetSession(ctx context.Context, userEmail string, refreshToken string, passHash []byte) error {

	setRef := r.Client.HSet(userEmail, refreshField, refreshToken)
	if setRef.Err() != nil {
		return setRef.Err()
	}

	setPass := r.Client.HSet(userEmail, passField, passHash)
	if setPass.Err() != nil {
		return setPass.Err()
	}

	exp := r.Client.Expire(userEmail, r.RefreshTTL)
	if exp.Err() != nil {
		return exp.Err()

	}
	return nil
}

// UserRefresh returns refresh token
func (r *Casher) UserRefresh(ctx context.Context, userEmail string) (string, error) {
	var result string

	err := r.Client.HGet(userEmail, refreshField).Scan(&result)

	if err != nil {
		return "", err
	}
	return result, nil

}

// UserPassHash returns pass hash
func (r *Casher) UserPassHash(ctx context.Context, userEmail string) (string, error) {
	var result string

	err := r.Client.HGet(userEmail, passField).Scan(&result)

	if err != nil {
		return "", err
	}
	return result, nil

}

// BlockSession deletes user from casher
func (r *Casher) BlockSession(ctx context.Context, userEmail string) error {

	err := r.Client.HDel(userEmail, refreshField, passField).Err()

	if err != nil {
		return err
	}
	return nil
}

// UserSession returns refresh token and pass hash
func (r *Casher) UserSession(ctx context.Context, userEmail string) (string, string, error) {
	result, err := r.Client.HGetAll(userEmail).Result()
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
