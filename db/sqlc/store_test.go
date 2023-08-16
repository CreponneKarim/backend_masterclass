package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := newStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	
	//	run n concurent transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i:=0; i<n; i++ {
		go func() {
			ctx := context.Background()

			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				amount: amount,
			})

			errs <- err
			results <- result
		}()
	}


	existed := make(map[int]bool)
	//	check results
	for i:=0; i<n; i++ {
		err := <- errs
		require.NoError(t,err)

		result := <- results
		require.NotEmpty(t,result)

		//	check all result parameters
		transfer := result.Transfer
		require.NotEmpty(t,transfer)
		require.Equal(t,transfer.ToAccountID, account2.ID)
		require.Equal(t,transfer.FromAccountID, account1.ID)
		require.Equal(t,transfer.Amount, amount)
		require.NotZero(t,transfer.ID)
		require.NotZero(t,transfer.CreatedAt)


		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t,err)

		//	check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t,fromEntry)
		require.Equal(t,fromEntry.AccountID, account1.ID)
		require.Equal(t,fromEntry.Amount, -amount)
		require.NotZero(t,fromEntry.ID)
		require.NotZero(t,fromEntry.CreatedAt)

		_,err = store.GetEntry(context.Background(),fromEntry.ID)
		require.NoError(t,err)

		ToEntry := result.ToEntry
		require.NotEmpty(t,ToEntry)
		require.Equal(t,ToEntry.AccountID, account2.ID)
		require.Equal(t,ToEntry.Amount, amount)
		require.NotZero(t,ToEntry.ID)
		require.NotZero(t,ToEntry.CreatedAt)

		_,err = store.GetEntry(context.Background(),ToEntry.ID)
		require.NoError(t,err)

		//	TODO : check account balances here

		//	check the accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t,fromAccount)
		require.Equal(t,fromAccount.ID,account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t,toAccount)
		require.Equal(t,toAccount.ID,account2.ID)

		//	check balances


		diff1 := int64(account1.Balance- fromAccount.Balance)
		diff2 := int64(toAccount.Balance - account2.Balance)

		require.Equal(t,diff1,diff2)

		require.True(t,diff1>0)
		require.True(t,diff1%amount ==0)

		k := int(diff1 / amount)

		require.True(t, k>=1 && k <=n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//	check the final updated balance
	updatedAccount1,err := testQueries.GetAccount(context.Background(),account1.ID)
	require.NoError(t,err)
	
	updatedAccount2,err := testQueries.GetAccount(context.Background(),account2.ID)
	require.NoError(t,err)

	require.Equal(t,account1.Balance - int64(n) * amount, updatedAccount1.Balance )
	require.Equal(t,account2.Balance + int64(n) * amount, updatedAccount2.Balance )
}

func TestTransferTxDeadlock(t *testing.T) {
	store := newStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	
	//	run n concurent transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i:=0; i<n; i++ {
		
		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i % 2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID	
		}
		
		go func() {
			ctx := context.Background()
			

			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID: toAccountId,
				amount: amount,
			})

			errs <- err
		}()
	}


	//	check results
	for i:=0; i<n; i++ {
		err := <- errs
		require.NoError(t,err)
	}

	//	check the final updated balance
	updatedAccount1,err := testQueries.GetAccount(context.Background(),account1.ID)
	require.NoError(t,err)
	
	updatedAccount2,err := testQueries.GetAccount(context.Background(),account2.ID)
	require.NoError(t,err)

	require.Equal(t,account1.Balance , updatedAccount1.Balance )
	require.Equal(t,account2.Balance , updatedAccount2.Balance )
}