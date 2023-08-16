package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func newStore(db *sql.DB) *Store{
	return &Store{
		db: db,
		Queries: New(db),
	}
}

//	execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err!=nil {
		return err
	}
	q := New(tx)
	
	err = fn(q)
	if err!=nil {
		if rbErr := tx.Rollback(); err!=nil {
			return fmt.Errorf("tx err : %v, rb err : %v",err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID	int64 `json:"from_account_id"`
	ToAccountID		int64 `json:"to_account_id"`
	amount			int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer		Transfer 	`json:"transfer"`
	FromAccount		Account		`json:"from_account"`
	ToAccount		Account 	`json:"to_account"`
	FromEntry		Entry		`json:"from_entry"`
	ToEntry			Entry		`json:"to_entry"`
}


func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx,func (q *Queries) error{
		var err error

		result.Transfer, err = q.CreateTransfer(ctx,CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.amount,
		})

		if err!=nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.amount,
		})
		if err!=nil {
			return err
		}

		result.ToEntry,err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.amount,
		})
		if err!=nil {
			return err
		}

		//	TO-DO: UPDATE ACCOUNTS BALANCE
		if arg.FromAccountID<arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.amount, arg.ToAccountID, arg.amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.amount, arg.FromAccountID, -arg.amount)
		}
		return err
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountId1 int64,
	amount1 int64,
	accountId2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		ID: accountId1,
		Amount: amount1,
	})
	if err!=nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		ID: accountId2,
		Amount: amount2,
	})
	
	return
}