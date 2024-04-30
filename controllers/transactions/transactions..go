package transactions

import (
	"context"
	"net/http"

	accounts_dao "github.com/ashwin-m/transactions/daos/accounts"
	transactions_dao "github.com/ashwin-m/transactions/daos/transactions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type createTransactionRequest struct {
	SourceAccountId      int64   `json:"source_account_id"`
	DestinationAccountId int64   `json:"destination_account_id"`
	Amount               float64 `json:"amount"`
}

type handler struct {
	dbPool          *pgxpool.Pool
	accountsDao     accounts_dao.Dao
	transactionsDao transactions_dao.Dao
}

type Handler interface {
	RouteGroup(*gin.Engine)
}

func NewHandler(dbPool *pgxpool.Pool, accountsDao accounts_dao.Dao, transactionsDao transactions_dao.Dao) Handler {
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
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	sourceAccount, err := h.accountsDao.GetById(request.SourceAccountId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	sourceAccountBalance := sourceAccount.GetBalance()

	destinationAccount, err := h.accountsDao.GetById(request.DestinationAccountId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	destinationAccountBalance := destinationAccount.GetBalance()

	if sourceAccountBalance < request.Amount {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	txn, err := h.dbPool.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	err = h.transactionsDao.Create(txn, sourceAccount.GetId(), destinationAccount.GetId(), request.Amount)
	if err != nil {
		txn.Rollback(context.Background())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	newSourceAccountBalance := sourceAccountBalance - request.Amount
	_, err = h.accountsDao.UpdateBalance(txn, request.SourceAccountId, newSourceAccountBalance)
	if err != nil {
		txn.Rollback(context.Background())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	newDestinationAccountBalance := destinationAccountBalance + request.Amount
	_, err = h.accountsDao.UpdateBalance(txn, request.DestinationAccountId, newDestinationAccountBalance)
	if err != nil {
		txn.Rollback(context.Background())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	txn.Commit(context.Background())

}
