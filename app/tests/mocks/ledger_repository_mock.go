// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/pagination"
	"sync"
	"time"
)

// Ensure, that RepositoryMock does implement domain.Repository.
// If this is not the case, regenerate this file with moq.
var _ domain.Repository = &RepositoryMock{}

// RepositoryMock is a mock implementation of domain.Repository.
//
// 	func TestSomethingThatUsesRepository(t *testing.T) {
//
// 		// make and configure a mocked domain.Repository
// 		mockedRepository := &RepositoryMock{
// 			CreateTransactionFunc: func(contextMoqParam context.Context, transaction entities.Transaction) error {
// 				panic("mock out the CreateTransaction method")
// 			},
// 			GetAnalyticAccountBalanceFunc: func(contextMoqParam context.Context, account vos.Account) (vos.AccountBalance, error) {
// 				panic("mock out the GetAnalyticAccountBalance method")
// 			},
// 			GetBoundedAccountBalanceFunc: func(contextMoqParam context.Context, account vos.Account, timeMoqParam1 time.Time, timeMoqParam2 time.Time) (vos.AccountBalance, error) {
// 				panic("mock out the GetBoundedAccountBalance method")
// 			},
// 			GetSyntheticAccountBalanceFunc: func(contextMoqParam context.Context, account vos.Account) (vos.AccountBalance, error) {
// 				panic("mock out the GetSyntheticAccountBalance method")
// 			},
// 			GetSyntheticReportFunc: func(contextMoqParam context.Context, account vos.Account, n int, timeMoqParam1 time.Time, timeMoqParam2 time.Time) (*vos.SyntheticReport, error) {
// 				panic("mock out the GetSyntheticReport method")
// 			},
// 			ListAccountEntriesFunc: func(contextMoqParam context.Context, accountEntryRequest vos.AccountEntryRequest) ([]vos.AccountEntry, pagination.Cursor, error) {
// 				panic("mock out the ListAccountEntries method")
// 			},
// 		}
//
// 		// use mockedRepository in code that requires domain.Repository
// 		// and then make assertions.
//
// 	}
type RepositoryMock struct {
	// CreateTransactionFunc mocks the CreateTransaction method.
	CreateTransactionFunc func(contextMoqParam context.Context, transaction entities.Transaction) error

	// GetAnalyticAccountBalanceFunc mocks the GetAnalyticAccountBalance method.
	GetAnalyticAccountBalanceFunc func(contextMoqParam context.Context, account vos.Account) (vos.AccountBalance, error)

	// GetBoundedAccountBalanceFunc mocks the GetBoundedAccountBalance method.
	GetBoundedAccountBalanceFunc func(contextMoqParam context.Context, account vos.Account, timeMoqParam1 time.Time, timeMoqParam2 time.Time) (vos.AccountBalance, error)

	// GetSyntheticAccountBalanceFunc mocks the GetSyntheticAccountBalance method.
	GetSyntheticAccountBalanceFunc func(contextMoqParam context.Context, account vos.Account) (vos.AccountBalance, error)

	// GetSyntheticReportFunc mocks the GetSyntheticReport method.
	GetSyntheticReportFunc func(contextMoqParam context.Context, account vos.Account, n int, timeMoqParam1 time.Time, timeMoqParam2 time.Time) (*vos.SyntheticReport, error)

	// ListAccountEntriesFunc mocks the ListAccountEntries method.
	ListAccountEntriesFunc func(contextMoqParam context.Context, accountEntryRequest vos.AccountEntryRequest) ([]vos.AccountEntry, pagination.Cursor, error)

	// calls tracks calls to the methods.
	calls struct {
		// CreateTransaction holds details about calls to the CreateTransaction method.
		CreateTransaction []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Transaction is the transaction argument value.
			Transaction entities.Transaction
		}
		// GetAnalyticAccountBalance holds details about calls to the GetAnalyticAccountBalance method.
		GetAnalyticAccountBalance []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Account is the account argument value.
			Account vos.Account
		}
		// GetBoundedAccountBalance holds details about calls to the GetBoundedAccountBalance method.
		GetBoundedAccountBalance []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Account is the account argument value.
			Account vos.Account
			// TimeMoqParam1 is the timeMoqParam1 argument value.
			TimeMoqParam1 time.Time
			// TimeMoqParam2 is the timeMoqParam2 argument value.
			TimeMoqParam2 time.Time
		}
		// GetSyntheticAccountBalance holds details about calls to the GetSyntheticAccountBalance method.
		GetSyntheticAccountBalance []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Account is the account argument value.
			Account vos.Account
		}
		// GetSyntheticReport holds details about calls to the GetSyntheticReport method.
		GetSyntheticReport []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Account is the account argument value.
			Account vos.Account
			// N is the n argument value.
			N int
			// TimeMoqParam1 is the timeMoqParam1 argument value.
			TimeMoqParam1 time.Time
			// TimeMoqParam2 is the timeMoqParam2 argument value.
			TimeMoqParam2 time.Time
		}
		// ListAccountEntries holds details about calls to the ListAccountEntries method.
		ListAccountEntries []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// AccountEntryRequest is the accountEntryRequest argument value.
			AccountEntryRequest vos.AccountEntryRequest
		}
	}
	lockCreateTransaction          sync.RWMutex
	lockGetAnalyticAccountBalance  sync.RWMutex
	lockGetBoundedAccountBalance   sync.RWMutex
	lockGetSyntheticAccountBalance sync.RWMutex
	lockGetSyntheticReport         sync.RWMutex
	lockListAccountEntries         sync.RWMutex
}

// CreateTransaction calls CreateTransactionFunc.
func (mock *RepositoryMock) CreateTransaction(contextMoqParam context.Context, transaction entities.Transaction) error {
	if mock.CreateTransactionFunc == nil {
		panic("RepositoryMock.CreateTransactionFunc: method is nil but Repository.CreateTransaction was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Transaction     entities.Transaction
	}{
		ContextMoqParam: contextMoqParam,
		Transaction:     transaction,
	}
	mock.lockCreateTransaction.Lock()
	mock.calls.CreateTransaction = append(mock.calls.CreateTransaction, callInfo)
	mock.lockCreateTransaction.Unlock()
	return mock.CreateTransactionFunc(contextMoqParam, transaction)
}

// CreateTransactionCalls gets all the calls that were made to CreateTransaction.
// Check the length with:
//     len(mockedRepository.CreateTransactionCalls())
func (mock *RepositoryMock) CreateTransactionCalls() []struct {
	ContextMoqParam context.Context
	Transaction     entities.Transaction
} {
	var calls []struct {
		ContextMoqParam context.Context
		Transaction     entities.Transaction
	}
	mock.lockCreateTransaction.RLock()
	calls = mock.calls.CreateTransaction
	mock.lockCreateTransaction.RUnlock()
	return calls
}

// GetAnalyticAccountBalance calls GetAnalyticAccountBalanceFunc.
func (mock *RepositoryMock) GetAnalyticAccountBalance(contextMoqParam context.Context, account vos.Account) (vos.AccountBalance, error) {
	if mock.GetAnalyticAccountBalanceFunc == nil {
		panic("RepositoryMock.GetAnalyticAccountBalanceFunc: method is nil but Repository.GetAnalyticAccountBalance was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Account         vos.Account
	}{
		ContextMoqParam: contextMoqParam,
		Account:         account,
	}
	mock.lockGetAnalyticAccountBalance.Lock()
	mock.calls.GetAnalyticAccountBalance = append(mock.calls.GetAnalyticAccountBalance, callInfo)
	mock.lockGetAnalyticAccountBalance.Unlock()
	return mock.GetAnalyticAccountBalanceFunc(contextMoqParam, account)
}

// GetAnalyticAccountBalanceCalls gets all the calls that were made to GetAnalyticAccountBalance.
// Check the length with:
//     len(mockedRepository.GetAnalyticAccountBalanceCalls())
func (mock *RepositoryMock) GetAnalyticAccountBalanceCalls() []struct {
	ContextMoqParam context.Context
	Account         vos.Account
} {
	var calls []struct {
		ContextMoqParam context.Context
		Account         vos.Account
	}
	mock.lockGetAnalyticAccountBalance.RLock()
	calls = mock.calls.GetAnalyticAccountBalance
	mock.lockGetAnalyticAccountBalance.RUnlock()
	return calls
}

// GetBoundedAccountBalance calls GetBoundedAccountBalanceFunc.
func (mock *RepositoryMock) GetBoundedAccountBalance(contextMoqParam context.Context, account vos.Account, timeMoqParam1 time.Time, timeMoqParam2 time.Time) (vos.AccountBalance, error) {
	if mock.GetBoundedAccountBalanceFunc == nil {
		panic("RepositoryMock.GetBoundedAccountBalanceFunc: method is nil but Repository.GetBoundedAccountBalance was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Account         vos.Account
		TimeMoqParam1   time.Time
		TimeMoqParam2   time.Time
	}{
		ContextMoqParam: contextMoqParam,
		Account:         account,
		TimeMoqParam1:   timeMoqParam1,
		TimeMoqParam2:   timeMoqParam2,
	}
	mock.lockGetBoundedAccountBalance.Lock()
	mock.calls.GetBoundedAccountBalance = append(mock.calls.GetBoundedAccountBalance, callInfo)
	mock.lockGetBoundedAccountBalance.Unlock()
	return mock.GetBoundedAccountBalanceFunc(contextMoqParam, account, timeMoqParam1, timeMoqParam2)
}

// GetBoundedAccountBalanceCalls gets all the calls that were made to GetBoundedAccountBalance.
// Check the length with:
//     len(mockedRepository.GetBoundedAccountBalanceCalls())
func (mock *RepositoryMock) GetBoundedAccountBalanceCalls() []struct {
	ContextMoqParam context.Context
	Account         vos.Account
	TimeMoqParam1   time.Time
	TimeMoqParam2   time.Time
} {
	var calls []struct {
		ContextMoqParam context.Context
		Account         vos.Account
		TimeMoqParam1   time.Time
		TimeMoqParam2   time.Time
	}
	mock.lockGetBoundedAccountBalance.RLock()
	calls = mock.calls.GetBoundedAccountBalance
	mock.lockGetBoundedAccountBalance.RUnlock()
	return calls
}

// GetSyntheticAccountBalance calls GetSyntheticAccountBalanceFunc.
func (mock *RepositoryMock) GetSyntheticAccountBalance(contextMoqParam context.Context, account vos.Account) (vos.AccountBalance, error) {
	if mock.GetSyntheticAccountBalanceFunc == nil {
		panic("RepositoryMock.GetSyntheticAccountBalanceFunc: method is nil but Repository.GetSyntheticAccountBalance was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Account         vos.Account
	}{
		ContextMoqParam: contextMoqParam,
		Account:         account,
	}
	mock.lockGetSyntheticAccountBalance.Lock()
	mock.calls.GetSyntheticAccountBalance = append(mock.calls.GetSyntheticAccountBalance, callInfo)
	mock.lockGetSyntheticAccountBalance.Unlock()
	return mock.GetSyntheticAccountBalanceFunc(contextMoqParam, account)
}

// GetSyntheticAccountBalanceCalls gets all the calls that were made to GetSyntheticAccountBalance.
// Check the length with:
//     len(mockedRepository.GetSyntheticAccountBalanceCalls())
func (mock *RepositoryMock) GetSyntheticAccountBalanceCalls() []struct {
	ContextMoqParam context.Context
	Account         vos.Account
} {
	var calls []struct {
		ContextMoqParam context.Context
		Account         vos.Account
	}
	mock.lockGetSyntheticAccountBalance.RLock()
	calls = mock.calls.GetSyntheticAccountBalance
	mock.lockGetSyntheticAccountBalance.RUnlock()
	return calls
}

// GetSyntheticReport calls GetSyntheticReportFunc.
func (mock *RepositoryMock) GetSyntheticReport(contextMoqParam context.Context, account vos.Account, n int, timeMoqParam1 time.Time, timeMoqParam2 time.Time) (*vos.SyntheticReport, error) {
	if mock.GetSyntheticReportFunc == nil {
		panic("RepositoryMock.GetSyntheticReportFunc: method is nil but Repository.GetSyntheticReport was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Account         vos.Account
		N               int
		TimeMoqParam1   time.Time
		TimeMoqParam2   time.Time
	}{
		ContextMoqParam: contextMoqParam,
		Account:         account,
		N:               n,
		TimeMoqParam1:   timeMoqParam1,
		TimeMoqParam2:   timeMoqParam2,
	}
	mock.lockGetSyntheticReport.Lock()
	mock.calls.GetSyntheticReport = append(mock.calls.GetSyntheticReport, callInfo)
	mock.lockGetSyntheticReport.Unlock()
	return mock.GetSyntheticReportFunc(contextMoqParam, account, n, timeMoqParam1, timeMoqParam2)
}

// GetSyntheticReportCalls gets all the calls that were made to GetSyntheticReport.
// Check the length with:
//     len(mockedRepository.GetSyntheticReportCalls())
func (mock *RepositoryMock) GetSyntheticReportCalls() []struct {
	ContextMoqParam context.Context
	Account         vos.Account
	N               int
	TimeMoqParam1   time.Time
	TimeMoqParam2   time.Time
} {
	var calls []struct {
		ContextMoqParam context.Context
		Account         vos.Account
		N               int
		TimeMoqParam1   time.Time
		TimeMoqParam2   time.Time
	}
	mock.lockGetSyntheticReport.RLock()
	calls = mock.calls.GetSyntheticReport
	mock.lockGetSyntheticReport.RUnlock()
	return calls
}

// ListAccountEntries calls ListAccountEntriesFunc.
func (mock *RepositoryMock) ListAccountEntries(contextMoqParam context.Context, accountEntryRequest vos.AccountEntryRequest) ([]vos.AccountEntry, pagination.Cursor, error) {
	if mock.ListAccountEntriesFunc == nil {
		panic("RepositoryMock.ListAccountEntriesFunc: method is nil but Repository.ListAccountEntries was just called")
	}
	callInfo := struct {
		ContextMoqParam     context.Context
		AccountEntryRequest vos.AccountEntryRequest
	}{
		ContextMoqParam:     contextMoqParam,
		AccountEntryRequest: accountEntryRequest,
	}
	mock.lockListAccountEntries.Lock()
	mock.calls.ListAccountEntries = append(mock.calls.ListAccountEntries, callInfo)
	mock.lockListAccountEntries.Unlock()
	return mock.ListAccountEntriesFunc(contextMoqParam, accountEntryRequest)
}

// ListAccountEntriesCalls gets all the calls that were made to ListAccountEntries.
// Check the length with:
//     len(mockedRepository.ListAccountEntriesCalls())
func (mock *RepositoryMock) ListAccountEntriesCalls() []struct {
	ContextMoqParam     context.Context
	AccountEntryRequest vos.AccountEntryRequest
} {
	var calls []struct {
		ContextMoqParam     context.Context
		AccountEntryRequest vos.AccountEntryRequest
	}
	mock.lockListAccountEntries.RLock()
	calls = mock.calls.ListAccountEntries
	mock.lockListAccountEntries.RUnlock()
	return calls
}
