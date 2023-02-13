package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < 5; i++ {
		go func() {
			res, err := store.TransferTx(context.Background(), TransferParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- res
		}()
	}

	for i := 0; i < 5; i++ {
		err := <-errs
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res)

		transfer := res.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, amount, transfer.Amount)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, acc1.ID)

		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, acc2.ID)

		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % amount == 0)
	}

	updateAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updateAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance - int64(5) * amount, updateAcc1.Balance)
	require.Equal(t, acc2.Balance + int64(5) * amount, updateAcc2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < 10; i++ {
		fromAccountID := acc1.ID
		toAccountID := acc2.ID


		if i%2 == 1{
			fromAccountID = acc2.ID
			toAccountID = acc1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < 10; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updateAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updateAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance, updateAcc1.Balance)
	require.Equal(t, acc2.Balance, updateAcc2.Balance)
}