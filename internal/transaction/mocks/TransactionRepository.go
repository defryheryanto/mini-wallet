// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	transaction "github.com/defryheryanto/mini-wallet/internal/transaction"
	mock "github.com/stretchr/testify/mock"
)

// TransactionRepository is an autogenerated mock type for the TransactionRepository type
type TransactionRepository struct {
	mock.Mock
}

// FindById provides a mock function with given fields: ctx, id
func (_m *TransactionRepository) FindById(ctx context.Context, id string) (*transaction.Transaction, error) {
	ret := _m.Called(ctx, id)

	var r0 *transaction.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*transaction.Transaction, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *transaction.Transaction); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transaction.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByReferenceId provides a mock function with given fields: ctx, referenceNo, transactionType
func (_m *TransactionRepository) FindByReferenceId(ctx context.Context, referenceNo string, transactionType string) (*transaction.Transaction, error) {
	ret := _m.Called(ctx, referenceNo, transactionType)

	var r0 *transaction.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*transaction.Transaction, error)); ok {
		return rf(ctx, referenceNo, transactionType)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *transaction.Transaction); ok {
		r0 = rf(ctx, referenceNo, transactionType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transaction.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, referenceNo, transactionType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTransactionsByWalletId provides a mock function with given fields: ctx, walletId
func (_m *TransactionRepository) FindTransactionsByWalletId(ctx context.Context, walletId string) ([]*transaction.Transaction, error) {
	ret := _m.Called(ctx, walletId)

	var r0 []*transaction.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*transaction.Transaction, error)); ok {
		return rf(ctx, walletId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*transaction.Transaction); ok {
		r0 = rf(ctx, walletId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*transaction.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, walletId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: ctx, data
func (_m *TransactionRepository) Insert(ctx context.Context, data *transaction.Transaction) error {
	ret := _m.Called(ctx, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *transaction.Transaction) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, data
func (_m *TransactionRepository) Update(ctx context.Context, data *transaction.Transaction) error {
	ret := _m.Called(ctx, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *transaction.Transaction) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTransactionRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewTransactionRepository creates a new instance of TransactionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTransactionRepository(t mockConstructorTestingTNewTransactionRepository) *TransactionRepository {
	mock := &TransactionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
