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

type accounts struct {
	Id      int64   `json:"id"`
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
	var account accounts

	err := c.ShouldBindJSON(&account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	accountModel, err := h.dao.Create(account.Id, account.Balance)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok && pgerrcode.IsIntegrityConstraintViolation(err.Code) {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	accountResponse := accounts{
		Id:      accountModel.GetId(),
		Balance: accountModel.GetBalance(),
	}

	c.JSON(http.StatusOK, accountResponse)

}

func (h *handler) get(c *gin.Context) {
	idString := c.Param("id")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	account, err := h.dao.GetById(id)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			c.JSON(http.StatusNotFound, gin.H{})
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
