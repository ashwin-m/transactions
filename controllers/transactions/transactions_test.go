package transactions

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	accountsdaomocks "github.com/ashwin-m/transactions/daos/accounts/mocks"
	transactionsdaomocks "github.com/ashwin-m/transactions/daos/transactions/mocks"
	accountsmodel "github.com/ashwin-m/transactions/models/accounts"
	"github.com/gin-gonic/gin"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionsCreate_BalancePassedAsInt(t *testing.T) {
	router := gin.Default()

	mockAccountsDao := accountsdaomocks.NewDao(t)
	mocktransactionsDao := transactionsdaomocks.NewDao(t)
	mockDB, _ := pgxmock.NewPool()

	h := NewHandler(mockDB, mockAccountsDao, mocktransactionsDao)
	h.RouteGroup(router)

	body := `{
		"source_account_id": 123,
		"destination_account_id": 456,
		"amount": 100
	}`
	bodyReader := strings.NewReader(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transactions", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"error\":\"json: cannot unmarshal number into Go struct field createTransactionRequest.amount of type string\"}", w.Body.String())
}

func TestTransactionsCreate_BalancePassedAsBadFloat(t *testing.T) {
	router := gin.Default()

	mockAccountsDao := accountsdaomocks.NewDao(t)
	mocktransactionsDao := transactionsdaomocks.NewDao(t)
	mockDB, _ := pgxmock.NewPool()

	h := NewHandler(mockDB, mockAccountsDao, mocktransactionsDao)
	h.RouteGroup(router)

	body := `{
		"source_account_id": 123,
		"destination_account_id": 456,
		"amount": "BC"
	}`
	bodyReader := strings.NewReader(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transactions", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"error\":\"unable to parse request amount\"}", w.Body.String())
}

func TestTransactionsCreate_SourceAccountDaoReturnsError(t *testing.T) {
	router := gin.Default()

	mockAccountsDao := accountsdaomocks.NewDao(t)
	mockAccountsDao.EXPECT().GetById(int64(123)).Return(accountsmodel.Accounts{}, errors.New("test"))
	mocktransactionsDao := transactionsdaomocks.NewDao(t)
	mockDB, _ := pgxmock.NewPool()

	h := NewHandler(mockDB, mockAccountsDao, mocktransactionsDao)
	h.RouteGroup(router)

	body := `{
		"source_account_id": 123,
		"destination_account_id": 456,
		"amount": "100.12345"
	}`
	bodyReader := strings.NewReader(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transactions", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{\"error\":\"test\"}", w.Body.String())
}

func TestTransactionsCreate_DestinationAccountDaoReturnsError(t *testing.T) {
	router := gin.Default()

	mockAccountsDao := accountsdaomocks.NewDao(t)

	sourceAccountId := int64(123)
	sourceAccount := accountsmodel.Accounts{}
	sourceAccount.SetId(sourceAccountId)
	sourceAccount.SetBalance(200.1)
	mockAccountsDao.EXPECT().GetById(sourceAccountId).Return(sourceAccount, nil)

	destinationAccountId := int64(456)
	destinationAccount := accountsmodel.Accounts{}
	mockAccountsDao.EXPECT().GetById(destinationAccountId).Return(destinationAccount, errors.New("test"))

	mocktransactionsDao := transactionsdaomocks.NewDao(t)
	mockDB, _ := pgxmock.NewPool()

	h := NewHandler(mockDB, mockAccountsDao, mocktransactionsDao)
	h.RouteGroup(router)

	body := `{
		"source_account_id": 123,
		"destination_account_id": 456,
		"amount": "100.12345"
	}`
	bodyReader := strings.NewReader(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transactions", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{\"error\":\"test\"}", w.Body.String())
}

func TestTransactionsCreate_SourceAccountHasLessBalance(t *testing.T) {
	router := gin.Default()

	mockAccountsDao := accountsdaomocks.NewDao(t)

	sourceAccountId := int64(123)
	sourceAccount := accountsmodel.Accounts{}
	sourceAccount.SetId(sourceAccountId)
	sourceAccount.SetBalance(100.1)
	mockAccountsDao.EXPECT().GetById(sourceAccountId).Return(sourceAccount, nil)

	mocktransactionsDao := transactionsdaomocks.NewDao(t)
	mockDB, _ := pgxmock.NewPool()

	h := NewHandler(mockDB, mockAccountsDao, mocktransactionsDao)
	h.RouteGroup(router)

	body := `{
		"source_account_id": 123,
		"destination_account_id": 456,
		"amount": "200.12345"
	}`
	bodyReader := strings.NewReader(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transactions", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"error\":\"account balance is less than transaction\"}", w.Body.String())
}

func TestTransactionsCreate_UnableToStartTxn(t *testing.T) {
	router := gin.Default()

	mockAccountsDao := accountsdaomocks.NewDao(t)

	sourceAccountId := int64(123)
	sourceAccount := accountsmodel.Accounts{}
	sourceAccount.SetId(sourceAccountId)
	sourceAccount.SetBalance(300.1)
	mockAccountsDao.EXPECT().GetById(sourceAccountId).Return(sourceAccount, nil)

	destinationAccountId := int64(456)
	destinationAccount := accountsmodel.Accounts{}
	destinationAccount.SetId(destinationAccountId)
	destinationAccount.SetBalance(200.1)
	mockAccountsDao.EXPECT().GetById(destinationAccountId).Return(destinationAccount, nil)

	mocktransactionsDao := transactionsdaomocks.NewDao(t)

	mockDB, _ := pgxmock.NewPool()
	mockDB.ExpectBegin().WillReturnError(errors.New("test"))

	h := NewHandler(mockDB, mockAccountsDao, mocktransactionsDao)
	h.RouteGroup(router)

	body := `{
		"source_account_id": 123,
		"destination_account_id": 456,
		"amount": "100.12345"
	}`
	bodyReader := strings.NewReader(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transactions", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"error\":\"test\"}", w.Body.String())
}

func TestTransactionsCreate_TransactionCreateReturnsError(t *testing.T) {
	router := gin.Default()

	amount := 100.12345

	mockAccountsDao := accountsdaomocks.NewDao(t)

	sourceAccountId := int64(123)
	sourceAccount := accountsmodel.Accounts{}
	sourceAccount.SetId(sourceAccountId)
	sourceAccount.SetBalance(300.1)
	mockAccountsDao.EXPECT().GetById(sourceAccountId).Return(sourceAccount, nil)

	destinationAccountId := int64(456)
	destinationAccount := accountsmodel.Accounts{}
	destinationAccount.SetId(destinationAccountId)
	destinationAccount.SetBalance(200.1)
	mockAccountsDao.EXPECT().GetById(destinationAccountId).Return(destinationAccount, nil)

	mocktransactionsDao := transactionsdaomocks.NewDao(t)
	mocktransactionsDao.EXPECT().Create(mock.Anything, sourceAccountId, destinationAccountId, amount).Return(0, errors.New("test"))

	mockDB, _ := pgxmock.NewPool()
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	h := NewHandler(mockDB, mockAccountsDao, mocktransactionsDao)
	h.RouteGroup(router)

	body := `{
		"source_account_id": 123,
		"destination_account_id": 456,
		"amount": "100.12345"
	}`
	bodyReader := strings.NewReader(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transactions", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"error\":\"test\"}", w.Body.String())
}

func TestTransactionsCreate_UpdateSourceAccountReturnsError(t *testing.T) {
	router := gin.Default()

	amount := 100.12345

	mockAccountsDao := accountsdaomocks.NewDao(t)

	sourceAccountId := int64(123)
	sourceAccount := accountsmodel.Accounts{}
	sourceAccount.SetId(sourceAccountId)
	sourceAccountBalance := 300.1
	sourceAccount.SetBalance(sourceAccountBalance)
	sourceVersion := int64(1)
	sourceAccount.SetVersion(sourceVersion)
	mockAccountsDao.EXPECT().GetById(sourceAccountId).Return(sourceAccount, nil)

	destinationAccountId := int64(456)
	destinationAccount := accountsmodel.Accounts{}
	destinationAccount.SetId(destinationAccountId)
	destinationAccountBalance := 200.1
	destinationAccount.SetBalance(destinationAccountBalance)
	mockAccountsDao.EXPECT().GetById(destinationAccountId).Return(destinationAccount, nil)

	mocktransactionsDao := transactionsdaomocks.NewDao(t)
	mocktransactionsDao.EXPECT().Create(mock.Anything, sourceAccountId, destinationAccountId, amount).Return(1, nil)

	newSourceAccountBalance := sourceAccountBalance - amount
	mockAccountsDao.EXPECT().UpdateBalance(mock.Anything, sourceAccountId, sourceVersion, newSourceAccountBalance).Return(sourceAccount, errors.New("test"))

	mockDB, _ := pgxmock.NewPool()
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	h := NewHandler(mockDB, mockAccountsDao, mocktransactionsDao)
	h.RouteGroup(router)

	body := `{
		"source_account_id": 123,
		"destination_account_id": 456,
		"amount": "100.12345"
	}`
	bodyReader := strings.NewReader(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transactions", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"error\":\"test\"}", w.Body.String())
}

func TestTransactionsCreate_UpdateDestinationAccountReturnsError(t *testing.T) {
	router := gin.Default()

	amount := 100.12345

	mockAccountsDao := accountsdaomocks.NewDao(t)

	sourceAccountId := int64(123)
	sourceAccount := accountsmodel.Accounts{}
	sourceAccount.SetId(sourceAccountId)
	sourceAccountBalance := 300.1
	sourceAccount.SetBalance(sourceAccountBalance)
	sourceVersion := int64(1)
	sourceAccount.SetVersion(sourceVersion)
	mockAccountsDao.EXPECT().GetById(sourceAccountId).Return(sourceAccount, nil)

	destinationAccountId := int64(456)
	destinationAccount := accountsmodel.Accounts{}
	destinationAccount.SetId(destinationAccountId)
	destinationAccountBalance := 200.1
	destinationAccount.SetBalance(destinationAccountBalance)
	destinationVersion := int64(2)
	destinationAccount.SetVersion(destinationVersion)
	mockAccountsDao.EXPECT().GetById(destinationAccountId).Return(destinationAccount, nil)

	mocktransactionsDao := transactionsdaomocks.NewDao(t)
	mocktransactionsDao.EXPECT().Create(mock.Anything, sourceAccountId, destinationAccountId, amount).Return(1, nil)

	newSourceAccountBalance := sourceAccountBalance - amount
	mockAccountsDao.EXPECT().UpdateBalance(mock.Anything, sourceAccountId, sourceVersion, newSourceAccountBalance).Return(sourceAccount, nil)

	newDestinationAccountBalance := destinationAccountBalance + amount
	mockAccountsDao.EXPECT().UpdateBalance(mock.Anything, destinationAccountId, destinationVersion, newDestinationAccountBalance).Return(sourceAccount, errors.New("test"))

	mockDB, _ := pgxmock.NewPool()
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	h := NewHandler(mockDB, mockAccountsDao, mocktransactionsDao)
	h.RouteGroup(router)

	body := `{
		"source_account_id": 123,
		"destination_account_id": 456,
		"amount": "100.12345"
	}`
	bodyReader := strings.NewReader(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transactions", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"error\":\"test\"}", w.Body.String())
}

func TestTransactionsCreate_Success(t *testing.T) {
	router := gin.Default()

	amount := 100.12345

	mockAccountsDao := accountsdaomocks.NewDao(t)

	sourceAccountId := int64(123)
	sourceAccount := accountsmodel.Accounts{}
	sourceAccount.SetId(sourceAccountId)
	sourceAccountBalance := 300.1
	sourceAccount.SetBalance(sourceAccountBalance)
	sourceVersion := int64(1)
	sourceAccount.SetVersion(sourceVersion)
	mockAccountsDao.EXPECT().GetById(sourceAccountId).Return(sourceAccount, nil)

	destinationAccountId := int64(456)
	destinationAccount := accountsmodel.Accounts{}
	destinationAccount.SetId(destinationAccountId)
	destinationAccountBalance := 200.1
	destinationAccount.SetBalance(destinationAccountBalance)
	destinationVersion := int64(2)
	destinationAccount.SetVersion(destinationVersion)
	mockAccountsDao.EXPECT().GetById(destinationAccountId).Return(destinationAccount, nil)

	mocktransactionsDao := transactionsdaomocks.NewDao(t)
	mocktransactionsDao.EXPECT().Create(mock.Anything, sourceAccountId, destinationAccountId, amount).Return(1, nil)

	newSourceAccountBalance := sourceAccountBalance - amount
	mockAccountsDao.EXPECT().UpdateBalance(mock.Anything, sourceAccountId, sourceVersion, newSourceAccountBalance).Return(sourceAccount, nil)

	newDestinationAccountBalance := destinationAccountBalance + amount
	mockAccountsDao.EXPECT().UpdateBalance(mock.Anything, destinationAccountId, destinationVersion, newDestinationAccountBalance).Return(sourceAccount, nil)

	mockDB, _ := pgxmock.NewPool()
	mockDB.ExpectBegin()
	mockDB.ExpectCommit()

	h := NewHandler(mockDB, mockAccountsDao, mocktransactionsDao)
	h.RouteGroup(router)

	body := `{
		"source_account_id": 123,
		"destination_account_id": 456,
		"amount": "100.12345"
	}`
	bodyReader := strings.NewReader(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transactions", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"transaction_id\":1}", w.Body.String())
}
