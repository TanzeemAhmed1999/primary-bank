package db

import (
	"context"
	"testing"
	"time"

	commonutils "github.com/primarybank/common/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	user := CreateRandomUser(t)

	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  commonutils.RandomMoney(),
		Currency: commonutils.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.NotEmpty(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)

	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	args := UpdateAccountParams{
		ID:      account1.ID,
		Balance: commonutils.RandomMoney(),
	}

	err := testStore.UpdateAccount(context.Background(), args)
	require.NoError(t, err)

	updatedAcc, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account1)

	require.Equal(t, account1.ID, updatedAcc.ID)
	require.Equal(t, account1.Owner, updatedAcc.Owner)
	require.Equal(t, args.Balance, updatedAcc.Balance)
	require.Equal(t, account1.Currency, updatedAcc.Currency)
	require.WithinDuration(t, account1.CreatedAt, updatedAcc.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	err := testStore.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	deletedAcc, err := testStore.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.Empty(t, deletedAcc)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testStore.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, len(accounts), 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestAddAccountBalanceAndGetAccountForUpdate(t *testing.T) {
	account := CreateRandomAccount(t)

	updateAmount := int64(50)
	updatedAccount, err := testStore.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		Amount: updateAmount,
		ID:     account.ID,
	})
	require.NoError(t, err)
	require.Equal(t, account.Balance+updateAmount, updatedAccount.Balance)

	fetchedAccount, err := testStore.GetAccountForUpdate(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, updatedAccount.Balance, fetchedAccount.Balance)
	require.Equal(t, updatedAccount.Owner, fetchedAccount.Owner)
	require.Equal(t, updatedAccount.Currency, fetchedAccount.Currency)
}
