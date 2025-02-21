package token

import "errors"

var (
	ErrExpiredToken = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)
