package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/muling3/bank-go/util"
	"github.com/stretchr/testify/require"
)

func createTestTransfer(t *testing.T) (Transfer, Account) {
	acc := createRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: acc.ID,
		ToAccountID: acc.ID,
		Amount:      util.RandomAmount(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NotEmpty(t, transfer)
	require.NoError(t, err)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer, acc
}

func TestCreateTransfer(t *testing.T) {
	createTestTransfer(t)

}

func TestGetTransfer(t *testing.T) {
	transfer, _ := createTestTransfer(t)


	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NotEmpty(t, transfer2)
	require.NoError(t, err)

	require.Equal(t, transfer.ID, transfer2.ID)
	require.Equal(t, transfer2.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer2.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer.Amount, transfer2.Amount)
}

func TestUpdateTransfer(t *testing.T) {
	transfer, acc := createTestTransfer(t)
	arg := UpdateTransferParams{
		ID:      transfer.ID,
		Amount: util.RandomAmount(),
	}

	transfer2, err := testQueries.UpdateTransfer(context.Background(), arg)

	require.NotEmpty(t, transfer2)
	require.NoError(t, err)

	require.Equal(t, transfer.ID, transfer2.ID)
	require.Equal(t, acc.ID, transfer2.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, arg.Amount, transfer2.Amount)
}

func TestDeleteTransfer(t *testing.T) {
	transfer, _ := createTestTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}

func TestListTransfers(t *testing.T) {

	for i := 0; i < 10; i++ {
		createTestTransfer(t)
	}

	transfers, err := testQueries.ListTransfers(context.Background(), 5)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}

}