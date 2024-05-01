package accounts

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	daoMocks "github.com/ashwin-m/transactions/daos/accounts/mocks"
	accounts_model "github.com/ashwin-m/transactions/models/accounts"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestAccountsCreate_BalancePassedAsInt(t *testing.T) {
	router := gin.Default()

	mockDao := daoMocks.NewDao(t)

	h := NewHandler(mockDao)
	h.RouteGroup(router)

	body := `{
		"account_id": 123,
		"initial_balance": 123
	}`
	bodyReader := strings.NewReader(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/accounts", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"error\":\"json: cannot unmarshal number into Go struct field createAccountsRequest.initial_balance of type string\"}", w.Body.String())
}

func TestAccountsCreate_BadFloatPassed(t *testing.T) {
	router := gin.Default()

	mockDao := daoMocks.NewDao(t)

	h := NewHandler(mockDao)
	h.RouteGroup(router)

	body := `{
		"account_id": 123,
		"initial_balance": "abc"
	}`
	bodyReader := strings.NewReader(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/accounts", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"error\":\"strconv.ParseFloat: parsing \\\"abc\\\": invalid syntax\"}", w.Body.String())
}

func TestAccountsCreate_DaoReturnError(t *testing.T) {
	router := gin.Default()

	mockDao := daoMocks.NewDao(t)
	mockDao.EXPECT().Create(int64(123), 100.23344).Return(accounts_model.Accounts{}, errors.New("test"))

	h := NewHandler(mockDao)
	h.RouteGroup(router)

	body := `{
		"account_id": 123,
		"initial_balance": "100.23344"
	}`
	bodyReader := strings.NewReader(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/accounts", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"error\":\"test\"}", w.Body.String())
}

func TestAccountsCreate_Success(t *testing.T) {
	router := gin.Default()

	mockDao := daoMocks.NewDao(t)
	mockDao.EXPECT().Create(int64(123), 100.23344).Return(accounts_model.Accounts{}, nil)

	h := NewHandler(mockDao)
	h.RouteGroup(router)

	body := `{
		"account_id": 123,
		"initial_balance": "100.23344"
	}`
	bodyReader := strings.NewReader(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/accounts", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body)
}

func TestAccountsGet_BadAccountId(t *testing.T) {
	router := gin.Default()

	mockDao := daoMocks.NewDao(t)

	h := NewHandler(mockDao)
	h.RouteGroup(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/accounts/abc", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"error\":\"strconv.ParseInt: parsing \\\"abc\\\": invalid syntax\"}", w.Body.String())
}

func TestAccountsGet_DaoReturnNoRowsError(t *testing.T) {
	router := gin.Default()

	accountId := int64(123)

	mockDao := daoMocks.NewDao(t)
	mockDao.EXPECT().GetById(accountId).Return(accounts_model.Accounts{}, pgx.ErrNoRows)

	h := NewHandler(mockDao)
	h.RouteGroup(router)

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/accounts/%d", accountId)
	req, _ := http.NewRequest("GET", url, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{\"error\":\"no rows in result set\"}", w.Body.String())
}

func TestAccountsGet_DaoReturnError(t *testing.T) {
	router := gin.Default()

	accountId := int64(123)

	mockDao := daoMocks.NewDao(t)
	mockDao.EXPECT().GetById(accountId).Return(accounts_model.Accounts{}, errors.New("test"))

	h := NewHandler(mockDao)
	h.RouteGroup(router)

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/accounts/%d", accountId)
	req, _ := http.NewRequest("GET", url, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"error\":\"test\"}", w.Body.String())
}

func TestAccountsGet_Success(t *testing.T) {
	router := gin.Default()

	accountId := int64(123)
	balance := 123.234

	mockDao := daoMocks.NewDao(t)

	account := accounts_model.Accounts{}
	account.SetId(accountId)
	account.SetBalance(balance)
	mockDao.EXPECT().GetById(accountId).Return(account, nil)

	expectedResponse := "{\"account_id\":123,\"balance\":123.234}"

	h := NewHandler(mockDao)
	h.RouteGroup(router)

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/accounts/%d", accountId)
	req, _ := http.NewRequest("GET", url, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}
