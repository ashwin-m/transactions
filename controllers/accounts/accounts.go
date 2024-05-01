package accounts

import (
	"net/http"
	"strconv"

	accounts_dao "github.com/ashwin-m/transactions/daos/accounts"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type createAccountsRequest struct {
	Id      int64  `json:"account_id"`
	Balance string `json:"initial_balance"`
}

type accounts struct {
	Id      int64   `json:"account_id"`
	Balance float64 `json:"balance"`
}

type handler struct {
	dao accounts_dao.Dao
}

type Handler interface {
	RouteGroup(*gin.Engine)
}

func NewHandler(dao accounts_dao.Dao) Handler {
	return &handler{
		dao: dao,
	}
}

func (h *handler) RouteGroup(r *gin.Engine) {
	rg := r.Group("/accounts")

	rg.POST("", h.create)
	rg.GET("/:id", h.get)
}

func (h *handler) create(c *gin.Context) {
	var account createAccountsRequest

	err := c.ShouldBindJSON(&account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	initialAccountBalance, err := strconv.ParseFloat(account.Balance, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.dao.Create(account.Id, initialAccountBalance)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok && pgerrcode.IsIntegrityConstraintViolation(err.Code) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, "")

}

func (h *handler) get(c *gin.Context) {
	idString := c.Param("id")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.dao.GetById(id)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	accountResponse := accounts{
		Id:      account.GetId(),
		Balance: account.GetBalance(),
	}

	c.JSON(http.StatusOK, accountResponse)
}
