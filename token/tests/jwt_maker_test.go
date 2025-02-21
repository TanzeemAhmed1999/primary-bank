package tests

import (
	"testing"
	"time"

	"github.com/primarybank/token"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	secretKey := "a_very_secure_secret_key_with_min_length"
	maker, err := token.NewJWTMaker(secretKey)
	require.NoError(t, err)
	duration := time.Minute
	username := "test_user"
	now := time.Now()
	tokenStr, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	payload, err := maker.VerifyToken(tokenStr)
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, now.Add(duration), payload.ExpiresAt.Time, time.Second)
}

func TestExpiredToken(t *testing.T) {
	secretKey := "a_very_secure_secret_key_with_min_length"
	maker, err := token.NewJWTMaker(secretKey)
	require.NoError(t, err)

	tokenStr, err := maker.CreateToken("test_user", -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	payload, err := maker.VerifyToken(tokenStr)
	require.ErrorIs(t, err, token.ErrInvalidToken)
	require.Nil(t, payload)
}

func TestInvalidToken(t *testing.T) {
	secretKey := "a_very_secure_secret_key_with_min_length"
	maker, err := token.NewJWTMaker(secretKey)
	require.NoError(t, err)

	invalidToken := "invalid.token.string"
	payload, err := maker.VerifyToken(invalidToken)
	require.Error(t, err, token.ErrInvalidToken)
	require.Nil(t, payload)
}
