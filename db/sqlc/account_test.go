package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/CreponneKarim/backend_masterclass/util"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {

	arg := CreateAccountParams{
		Owner: util.RandomOwnerName(),
		Balance: util.RandomBalance(),
		Curency: util.RandomCurrency(),
	}	
	
	account,err := testQueries.CreateAccount(context.Background(), arg)
	
	require.NoError(t,err)
	require.NotEmpty(t,account)
	require.Equal(t,arg.Owner,account.Owner)
	require.Equal(t,arg.Balance,account.Balance)
	require.Equal(t,arg.Curency,account.Curency)
	require.NotZero(t,account.ID)
	require.NotZero(t,account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T){
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T){
	account1 := createRandomAccount(t)

	account2, err := testQueries.GetAccount(context.Background(),account1.ID)

	require.NoError(t,err)
	require.NotEmpty(t,account2)
	require.Equal(t,account1.ID,account2.ID)
	require.WithinDuration(t,account1.CreatedAt,account2.CreatedAt,time.Second)
	require.Equal(t,account1.Owner,account2.Owner)
	require.Equal(t,account1.Balance,account2.Balance)
	require.Equal(t,account1.Curency,account2.Curency)
}

func TestUpdateAccount(t *testing.T){
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams {
		ID: account1.ID,
		Balance: util.RandomBalance(),
	}

	account2,err := testQueries.UpdateAccount(context.Background(),arg)

	require.NoError(t,err)
	require.NotEmpty(t,account2)

	require.Equal(t,account1.Owner,account2.Owner)
	require.Equal(t,arg.Balance,account2.Balance)
	require.Equal(t,account1.Curency,account2.Curency)
	require.WithinDuration(t,account1.CreatedAt,account2.CreatedAt,time.Second)
}

func TestDeleteAccount(t *testing.T){
	account1 :=createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(),account1.ID)
	require.NoError(t,err)

	account2, err := testQueries.GetAccount(context.Background(),account1.ID)

	require.Error(t,err)
	require.Equal(t,err,sql.ErrNoRows)
	require.Empty(t,account2)
}

func TestListAccount(t *testing.T){
	for i:=0; i<10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams {
		Limit: 5,
		Offset: 5,
	}

	accounts, err :=testQueries.ListAccounts(context.Background(),arg)

	require.NoError(t,err)
	require.Len(t,accounts,5)

	for _,account := range accounts {
		require.NotEmpty(t,account)
	}
}