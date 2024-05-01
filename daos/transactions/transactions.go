package transactionledger

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockery --name=Dao --output=mocks --outpkg=mocks --with-expecter
type Dao interface {
	Create(txn pgx.Tx, sourceAccountId, destinationAccountId int64, amount float64) (int64, error)
}

type dao struct {
	dbPool *pgxpool.Pool
}

func NewDao(dbPool *pgxpool.Pool) Dao {
	return &dao{
		dbPool: dbPool,
	}
}

func (d *dao) Create(txn pgx.Tx, sourceAccountId, destinationAccountId int64, amount float64) (int64, error) {
	var transactionId int64
	sqlStatement := "insert into transactions(source_account_id, destination_account_id, amount) values ($1, $2, $3)"
	err := txn.QueryRow(context.Background(), sqlStatement, sourceAccountId, destinationAccountId, amount).Scan(&transactionId)

	return transactionId, err
}
