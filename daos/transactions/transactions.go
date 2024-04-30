package transactionledger

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockery --name=Dao --output=mocks --outpkg=mocks
type Dao interface {
	Create(txn pgx.Tx, sourceAccountId, destinationAccountId int64, amount float64) error
}

type dao struct {
	dbPool *pgxpool.Pool
}

func NewDao(dbPool *pgxpool.Pool) Dao {
	return &dao{
		dbPool: dbPool,
	}
}

func (d *dao) Create(txn pgx.Tx, sourceAccountId, destinationAccountId int64, amount float64) error {
	sqlStatement := "insert into transactions(source_account_id, destination_account_id, amount) values ($1, $2, $3)"
	_, err := txn.Exec(context.Background(), sqlStatement, sourceAccountId, destinationAccountId, amount)

	return err
}
