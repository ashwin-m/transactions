package transactions

import (
	"context"
	"errors"
	"math/big"
	"net/http"

	accountsdao "github.com/ashwin-m/transactions/daos/accounts"
	transactionsdao "github.com/ashwin-m/transactions/daos/transactions"
	accountsmodel "github.com/ashwin-m/transactions/models/accounts"
	"github.com/ashwin-m/transactions/utils/pgxiface"
	"github.com/gin-gonic/gin"
)

const (
	prec                                = 5
	min_transaction_amount              = 0
	min_account_balance_for_transaction = 0
)

type createTransactionRequest struct {
	SourceAccountId      int64  `json:"source_account_id"`
	DestinationAccountId int64  `json:"destination_account_id"`
	Amount               string `json:"amount"`
}

type handler struct {
	dbPool          pgxiface.PgxIface
	accountsDao     accountsdao.Dao
	transactionsDao transactionsdao.Dao
}

type Handler interface {
	RouteGroup(*gin.Engine)
}

func NewHandler(dbPool pgxiface.PgxIface, accountsDao accountsdao.Dao, transactionsDao transactionsdao.Dao) Handler {
	return &handler{
		dbPool:          dbPool,
		accountsDao:     accountsDao,
		transactionsDao: transactionsDao,
	}
}

func (h *handler) RouteGroup(r *gin.Engine) {
	rg := r.Group("/transactions")

	rg.POST("", h.create)
}

func (h *handler) create(c *gin.Context) {
	var request createTransactionRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount, ok := new(big.Float).SetPrec(prec).SetString(request.Amount)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to parse request amount"})
		return
	}

	if amount.Cmp(big.NewFloat(min_transaction_amount)) == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request amount cant be less than 0"})
		return
	}

	sourceAccount, err := h.accountsDao.GetById(request.SourceAccountId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	sourceAccountBalanceFloat := sourceAccount.GetBalance()
	sourceAccountBalance := new(big.Float).SetPrec(prec).SetFloat64(sourceAccountBalanceFloat)

	err = validateSourceAccount(sourceAccount, amount)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	destinationAccount, err := h.accountsDao.GetById(request.DestinationAccountId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	destinationAccountBalanceFloat := destinationAccount.GetBalance()
	destinationAccountBalance := new(big.Float).SetPrec(prec).SetFloat64(destinationAccountBalanceFloat)

	txn, err := h.dbPool.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	amountFloat, _ := amount.Float64()
	transactionId, err := h.transactionsDao.Create(txn, sourceAccount.GetId(), destinationAccount.GetId(), amountFloat)
	if err != nil {
		txn.Rollback(context.Background())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newSourceAccountBalance := big.NewFloat(0).Sub(sourceAccountBalance, amount)
	newSourceAccountBalanceFloat, _ := newSourceAccountBalance.Float64()
	_, err = h.accountsDao.UpdateBalance(txn, request.SourceAccountId, sourceAccount.GetVersion(), newSourceAccountBalanceFloat)
	if err != nil {
		txn.Rollback(context.Background())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newDestinationAccountBalance := big.NewFloat(0).Sub(destinationAccountBalance, amount)
	newDestinationAccountBalanceFloat, _ := newDestinationAccountBalance.Float64()
	_, err = h.accountsDao.UpdateBalance(txn, request.DestinationAccountId, destinationAccount.GetVersion(), newDestinationAccountBalanceFloat)
	if err != nil {
		txn.Rollback(context.Background())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	txn.Commit(context.Background())
	c.JSON(http.StatusOK, gin.H{"transaction_id": transactionId})

}

func validateSourceAccount(sourceAccount accountsmodel.Accounts, transactionAmount *big.Float) error {
	sourceAccountBalanceFloat := sourceAccount.GetBalance()
	sourceAccountBalance := new(big.Float).SetPrec(prec).SetFloat64(sourceAccountBalanceFloat)

	if sourceAccountBalance.Cmp(transactionAmount) == -1 {
		return errors.New("account balance is less than transaction")
	}

	if sourceAccountBalance.Cmp(big.NewFloat(min_account_balance_for_transaction)) == -1 {
		return errors.New("account balance is less than minimum amount for transactions")
	}

	return nil
}
