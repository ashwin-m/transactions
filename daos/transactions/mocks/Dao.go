// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	pgx "github.com/jackc/pgx/v5"
	mock "github.com/stretchr/testify/mock"
)

// Dao is an autogenerated mock type for the Dao type
type Dao struct {
	mock.Mock
}

type Dao_Expecter struct {
	mock *mock.Mock
}

func (_m *Dao) EXPECT() *Dao_Expecter {
	return &Dao_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: txn, sourceAccountId, destinationAccountId, amount
func (_m *Dao) Create(txn pgx.Tx, sourceAccountId int64, destinationAccountId int64, amount float64) error {
	ret := _m.Called(txn, sourceAccountId, destinationAccountId, amount)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(pgx.Tx, int64, int64, float64) error); ok {
		r0 = rf(txn, sourceAccountId, destinationAccountId, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Dao_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type Dao_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - txn pgx.Tx
//   - sourceAccountId int64
//   - destinationAccountId int64
//   - amount float64
func (_e *Dao_Expecter) Create(txn interface{}, sourceAccountId interface{}, destinationAccountId interface{}, amount interface{}) *Dao_Create_Call {
	return &Dao_Create_Call{Call: _e.mock.On("Create", txn, sourceAccountId, destinationAccountId, amount)}
}

func (_c *Dao_Create_Call) Run(run func(txn pgx.Tx, sourceAccountId int64, destinationAccountId int64, amount float64)) *Dao_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(pgx.Tx), args[1].(int64), args[2].(int64), args[3].(float64))
	})
	return _c
}

func (_c *Dao_Create_Call) Return(_a0 error) *Dao_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Dao_Create_Call) RunAndReturn(run func(pgx.Tx, int64, int64, float64) error) *Dao_Create_Call {
	_c.Call.Return(run)
	return _c
}

// NewDao creates a new instance of Dao. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDao(t interface {
	mock.TestingT
	Cleanup(func())
}) *Dao {
	mock := &Dao{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}