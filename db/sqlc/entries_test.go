package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/muling3/bank-go/util"
	"github.com/stretchr/testify/require"
)

func createTestEntry(t *testing.T) (Entry, Account) {
	acc := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    util.RandomAmount(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NotEmpty(t, entry)
	require.NoError(t, err)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry, acc
}

func TestCreateEntry(t *testing.T) {
	createTestEntry(t)

}

func TestGetEntry(t *testing.T) {
	entry, acc := createTestEntry(t)


	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.NotEmpty(t, entry2)
	require.NoError(t, err)

	require.Equal(t, entry.ID, entry2.ID)
	require.Equal(t, acc.ID, entry2.AccountID)
	require.Equal(t, entry.Amount, entry2.Amount)
	require.WithinDuration(t, entry.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	entry, acc := createTestEntry(t)
	arg := UpdateEntryParams{
		ID:      entry.ID,
		Amount: util.RandomAmount(),
	}

	entry2, err := testQueries.UpdateEntry(context.Background(), arg)

	require.NotEmpty(t, entry2)
	require.NoError(t, err)

	require.Equal(t, entry.ID, entry2.ID)
	require.Equal(t, entry.AccountID, entry2.AccountID)
	require.Equal(t, acc.ID, entry2.AccountID)
	require.Equal(t, arg.Amount, entry2.Amount)
	require.WithinDuration(t, entry.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	entry, _ := createTestEntry(t)
	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestListEntries(t *testing.T) {

	for i := 0; i < 10; i++ {
		createTestEntry(t)
	}

	entries, err := testQueries.ListEntries(context.Background(), 5)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}

}
