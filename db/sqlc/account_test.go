package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/muling3/bank-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), arg)

	require.NotEmpty(t, acc)
	require.NoError(t, err)

	require.Equal(t, arg.Owner, acc.Owner)
	require.Equal(t, arg.Balance, acc.Balance)
	require.Equal(t, arg.Currency, acc.Currency)

	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)

	return acc
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc := createRandomAccount(t)

	acc2, err := testQueries.GetAccount(context.Background(), acc.ID)

	require.NotEmpty(t, acc2)
	require.NoError(t, err)

	require.Equal(t, acc.ID, acc2.ID)
	require.Equal(t, acc.Owner, acc2.Owner)
	require.Equal(t, acc.Balance, acc2.Balance)
	require.Equal(t, acc.Currency, acc2.Currency)
	require.WithinDuration(t, acc.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	acc := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      acc.ID,
		Balance: util.RandomAmount(),
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NotEmpty(t, acc2)
	require.NoError(t, err)

	require.Equal(t, acc.ID, acc2.ID)
	require.Equal(t, acc.Owner, acc2.Owner)
	require.Equal(t, arg.Balance, acc2.Balance)
	require.Equal(t, acc.Currency, acc2.Currency)
	require.WithinDuration(t, acc.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acc := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), acc.ID)
	require.NoError(t, err)

	acc2, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc2)
}

func TestListAccounts(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	accs, err := testQueries.ListAccounts(context.Background(), 5)
	require.NoError(t, err)
	require.Len(t, accs, 5)

	for _, acc := range accs{
		require.NotEmpty(t, acc)
	}
	
}
