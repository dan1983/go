package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dan1983/go/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Owner, account.Owner)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)

}
func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	getAccount, err := testQueries.GetAccounts(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getAccount)

	require.Equal(t, account1.ID, getAccount.ID)
	require.Equal(t, account1.Owner, getAccount.Owner)
	require.Equal(t, account1.Balance, getAccount.Balance)
	require.Equal(t, account1.Currency, getAccount.Currency)

	require.WithinDuration(t, account1.CreatedAt, getAccount.CreatedAt, time.Second)

}
func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg2 := UpdateAccountParams{
		Balance: util.RandomMoney(),
		ID:      account1.ID,
	}

	err := testQueries.UpdateAccount(context.Background(), arg2)

	if err != nil {
		t.Errorf("error updating account: %s", err.Error())
		return
	}

	getUpdatedAccount, err := testQueries.GetAccounts(context.Background(), account1.ID)
	if err != nil {
		t.Errorf("error updating account: %s", err.Error())
		return
	}

	require.NotEmpty(t, getUpdatedAccount)

	require.Equal(t, account1.ID, getUpdatedAccount.ID)

	require.Equal(t, account1.Owner, getUpdatedAccount.Owner)

	require.Equal(t, arg2.Balance, getUpdatedAccount.Balance)

	require.WithinDuration(t, account1.CreatedAt, getUpdatedAccount.CreatedAt, time.Second)

}
func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	deletedAccount, err := testQueries.GetAccounts(context.Background(), account1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedAccount)

}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
