package db

import (
	"context"
	"testing"
	"time"

	commonutils "github.com/primarybank/common/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username: commonutils.RandomString(10),
		Password: commonutils.RandomString(16),
		FullName: commonutils.RandomString(10),
		Email:    commonutils.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	createdUser := CreateRandomUser(t)

	retrievedUser, err := testStore.GetUser(context.Background(), createdUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedUser)

	require.Equal(t, createdUser.Username, retrievedUser.Username)
	require.Equal(t, createdUser.Password, retrievedUser.Password)
	require.Equal(t, createdUser.FullName, retrievedUser.FullName)
	require.Equal(t, createdUser.Email, retrievedUser.Email)
	require.WithinDuration(t, createdUser.CreatedAt, retrievedUser.CreatedAt, time.Second)
	require.WithinDuration(t, createdUser.UpdatedAt, retrievedUser.UpdatedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	createdUser := CreateRandomUser(t)

	arg := UpdateUserParams{
		Username: createdUser.Username,
		FullName: commonutils.RandomString(12),
		Email:    commonutils.RandomEmail(),
	}

	updatedUser, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, createdUser.Username, updatedUser.Username)
	require.Equal(t, arg.FullName, updatedUser.FullName)
	require.Equal(t, arg.Email, updatedUser.Email)
	require.Equal(t, createdUser.Password, updatedUser.Password)

	require.WithinDuration(t, createdUser.CreatedAt, updatedUser.CreatedAt, time.Second)
	require.NotEqual(t, createdUser.UpdatedAt, updatedUser.UpdatedAt)
}
