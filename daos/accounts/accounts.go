package accounts

import (
	"context"

	accounts_model "github.com/ashwin-m/transactions/models/accounts"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dao interface {
	GetById(id int64) (accounts_model.Accounts, error)
	Create(id int64, balanace float64) (accounts_model.Accounts, error)
	UpdateBalance(tx pgx.Tx, id int64, newBalance float64) (accounts_model.Accounts, error)
}

type dao struct {
	dbPool *pgxpool.Pool
}

func NewDao(dbPool *pgxpool.Pool) Dao {
	return &dao{
		dbPool: dbPool,
	}
}

func (d *dao) GetById(id int64) (accounts_model.Accounts, error) {
	var accountId int64
	var balance float64
	var account accounts_model.Accounts

	sqlStatement := "select id, balance from Accounts where id=$1"
	err := d.dbPool.QueryRow(context.Background(), sqlStatement, id).Scan(&accountId, &balance)
	if err == nil {
		account.SetId(id)
		account.SetBalance(balance)
	}

	return account, err
}

func (d *dao) Create(id int64, balance float64) (accounts_model.Accounts, error) {
	var account accounts_model.Accounts
	sqlStatement := "insert into Accounts(id, balance) values ($1, $2)"
	_, err := d.dbPool.Exec(context.Background(), sqlStatement, id, balance)
	if err == nil {
		account.SetId(id)
		account.SetBalance(balance)
	}

	return account, err
}

func (d *dao) UpdateBalance(tx pgx.Tx, id int64, newBalance float64) (accounts_model.Accounts, error) {
	var account accounts_model.Accounts
	sqlStatement := "UPDATE accounts SET balance=$2 where id=$1"
	_, err := tx.Exec(context.Background(), sqlStatement, id, newBalance)
	if err == nil {
		account.SetId(id)
		account.SetBalance(newBalance)
	}

	return account, err
}
