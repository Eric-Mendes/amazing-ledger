package ledger

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/instrumentators"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/pagination"
)

func Test_generateListAccountEntriesQuery(t *testing.T) {
	t.Parallel()

	account, err := vos.NewAnalyticAccount("liability.test.account1")
	assert.NoError(t, err)

	synthAccount, err := vos.NewAccount("liability.*.account1")
	assert.NoError(t, err)

	size := 10

	end := time.Now().UTC().Round(time.Microsecond)
	start := end.Add(-10 * time.Second)

	version := vos.Version(1)

	testCases := []struct {
		name          string
		req           func() vos.AccountEntryRequest
		expectedQuery string
		expectedArgs  []interface{}
		expectedErr   error
	}{
		{
			name: "valid - no filters - no pagination",
			req: func() vos.AccountEntryRequest {
				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: start,
					EndDate:   end,
					Page: pagination.Page{
						Size:   size,
						Cursor: nil,
					},
				}
			},
			expectedQuery: fmt.Sprintf(_accountEntriesQueryPrefix, "=") + _accountEntriesQuerySuffixAnalytic,
			expectedArgs:  []interface{}{account.Value(), start, end, size + 1},
			expectedErr:   nil,
		},
		{
			name: "valid - no filters - no pagination",
			req: func() vos.AccountEntryRequest {
				return vos.AccountEntryRequest{
					Account:   synthAccount,
					StartDate: start,
					EndDate:   end,
					Page: pagination.Page{
						Size:   size,
						Cursor: nil,
					},
				}
			},
			expectedQuery: fmt.Sprintf(_accountEntriesQueryPrefix, "~") + _accountEntriesQuerySuffixSynthetic,
			expectedArgs:  []interface{}{synthAccount.Value(), start, end, size + 1},
			expectedErr:   nil,
		},
		{
			name: "valid - no filters - with pagination",
			req: func() vos.AccountEntryRequest {
				cursor, _ := pagination.NewCursor(listAccountEntriesCursor{
					CompetenceDate: end,
					Version:        1,
				})

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: start,
					EndDate:   end,
					Page: pagination.Page{
						Size:   size,
						Cursor: cursor,
					},
				}
			},
			expectedQuery: fmt.Sprintf(_accountEntriesQueryPrefix, "=") +
				fmt.Sprintf(_accountEntriesQueryPaginationAnalytic, 5, 6) +
				_accountEntriesQuerySuffixAnalytic,
			expectedArgs: []interface{}{account.Value(), start, end, size + 1, end, version.AsInt64()},
			expectedErr:  nil,
		},
		{
			name: "valid - with single company filter - no pagination",
			req: func() vos.AccountEntryRequest {
				filter := vos.AccountEntryFilter{Companies: []string{"company_1"}}

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: start,
					EndDate:   end,
					Filter:    filter,
					Page: pagination.Page{
						Size:   size,
						Cursor: nil,
					},
				}
			},
			expectedQuery: fmt.Sprintf(_accountEntriesQueryPrefix, "=") +
				fmt.Sprintf(_accountEntriesCompanyFilter, 5) +
				_accountEntriesQuerySuffixAnalytic,
			expectedArgs: []interface{}{account.Value(), start, end, size + 1, "company_1"},
			expectedErr:  nil,
		},
		{
			name: "valid - with multiple companies filter - no pagination",
			req: func() vos.AccountEntryRequest {
				filter := vos.AccountEntryFilter{Companies: []string{"company_1", "company_2"}}

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: start,
					EndDate:   end,
					Filter:    filter,
					Page: pagination.Page{
						Size:   size,
						Cursor: nil,
					},
				}
			},
			expectedQuery: fmt.Sprintf(_accountEntriesQueryPrefix, "=") +
				fmt.Sprintf(_accountEntriesCompaniesFilter, 5) +
				_accountEntriesQuerySuffixAnalytic,
			expectedArgs: []interface{}{account.Value(), start, end, size + 1, []string{"company_1", "company_2"}},
			expectedErr:  nil,
		},
		{
			name: "valid - with single event filter - no pagination",
			req: func() vos.AccountEntryRequest {
				filter := vos.AccountEntryFilter{Events: []int32{1}}

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: start,
					EndDate:   end,
					Filter:    filter,
					Page: pagination.Page{
						Size:   size,
						Cursor: nil,
					},
				}
			},
			expectedQuery: fmt.Sprintf(_accountEntriesQueryPrefix, "=") +
				fmt.Sprintf(_accountEntriesEventFilter, 5) +
				_accountEntriesQuerySuffixAnalytic,
			expectedArgs: []interface{}{account.Value(), start, end, size + 1, int32(1)},
			expectedErr:  nil,
		},
		{
			name: "valid - with multiple events filter - no pagination",
			req: func() vos.AccountEntryRequest {
				filter := vos.AccountEntryFilter{Events: []int32{1, 2}}

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: start,
					EndDate:   end,
					Filter:    filter,
					Page: pagination.Page{
						Size:   size,
						Cursor: nil,
					},
				}
			},
			expectedQuery: fmt.Sprintf(_accountEntriesQueryPrefix, "=") +
				fmt.Sprintf(_accountEntriesEventsFilter, 5) +
				_accountEntriesQuerySuffixAnalytic,
			expectedArgs: []interface{}{account.Value(), start, end, size + 1, []int32{1, 2}},
			expectedErr:  nil,
		},
		{
			name: "valid - with operation filter - no pagination",
			req: func() vos.AccountEntryRequest {
				filter := vos.AccountEntryFilter{Operation: vos.CreditOperation}

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: start,
					EndDate:   end,
					Filter:    filter,
					Page: pagination.Page{
						Size:   size,
						Cursor: nil,
					},
				}
			},
			expectedQuery: fmt.Sprintf(_accountEntriesQueryPrefix, "=") +
				fmt.Sprintf(_accountEntriesOperationFilter, 5) +
				_accountEntriesQuerySuffixAnalytic,
			expectedArgs: []interface{}{account.Value(), start, end, size + 1, vos.CreditOperation},
			expectedErr:  nil,
		},
		{
			name: "valid - all filters - with pagination",
			req: func() vos.AccountEntryRequest {
				filter := vos.AccountEntryFilter{
					Companies: []string{"company_1", "company_2"},
					Events:    []int32{1},
					Operation: vos.CreditOperation,
				}

				cursor, _ := pagination.NewCursor(listAccountEntriesCursor{
					CompetenceDate: end,
					Version:        1,
				})

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: start,
					EndDate:   end,
					Filter:    filter,
					Page: pagination.Page{
						Size:   size,
						Cursor: cursor,
					},
				}
			},
			expectedQuery: fmt.Sprintf(_accountEntriesQueryPrefix, "=") +
				fmt.Sprintf(_accountEntriesCompaniesFilter, 5) +
				fmt.Sprintf(_accountEntriesEventFilter, 6) +
				fmt.Sprintf(_accountEntriesOperationFilter, 7) +
				fmt.Sprintf(_accountEntriesQueryPaginationAnalytic, 8, 9) +
				_accountEntriesQuerySuffixAnalytic,
			expectedArgs: []interface{}{
				account.Value(), start, end, size + 1,
				[]string{"company_1", "company_2"}, int32(1), vos.CreditOperation,
				end, version.AsInt64(),
			},
			expectedErr: nil,
		},
		{
			name: "invalid page	token",
			req: func() vos.AccountEntryRequest {
				cursor, _ := pagination.NewCursor(map[string]interface{}{"version": "none"})

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: start,
					EndDate:   end,
					Page: pagination.Page{
						Size:   size,
						Cursor: cursor,
					},
				}
			},
			expectedQuery: "",
			expectedArgs:  nil,
			expectedErr:   app.ErrInvalidPageCursor,
		},
		{
			name: "invalid operation filter",
			req: func() vos.AccountEntryRequest {
				filter := vos.AccountEntryFilter{Operation: vos.InvalidOperation}

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: start,
					EndDate:   end,
					Filter:    filter,
					Page: pagination.Page{
						Size:   size,
						Cursor: nil,
					},
				}
			},
			expectedQuery: fmt.Sprintf(_accountEntriesQueryPrefix, "=") + _accountEntriesQuerySuffixAnalytic,
			expectedArgs:  []interface{}{account.Value(), start, end, size + 1},
			expectedErr:   nil,
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotQuery, gotArgs, err := generateListAccountEntriesQuery(tt.req())
			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, tt.expectedQuery, gotQuery)
			assert.EqualValues(t, tt.expectedArgs, gotArgs)
		})
	}
}

func TestLedgerRepository_ListAccountEntries(t *testing.T) {
	t.Parallel()

	type w struct {
		entries []vos.AccountEntry
		cursor  pagination.Cursor
	}

	const (
		account1     = "liability.abc.account1"
		account2     = "liability.abc.account2"
		synthAccount = "liability.abc.*"
		amount       = 100
	)

	testCases := []struct {
		name         string
		seedRepo     func(*testing.T, context.Context, *Repository) []entities.Transaction
		setupRequest func(*testing.T, []entities.Transaction, *Repository) vos.AccountEntryRequest
		want         func(*testing.T, []entities.Transaction, *Repository) w
	}{
		{
			name: "no exiting entries case",
			seedRepo: func(t *testing.T, ctx context.Context, r *Repository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), amount)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx := createTransaction(t, ctx, r, e1, e2)

				return []entities.Transaction{tx}
			},
			setupRequest: func(t *testing.T, _ []entities.Transaction, _ *Repository) vos.AccountEntryRequest {
				account, err := vos.NewAnalyticAccount("liability.abc.account3")
				assert.NoError(t, err)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Page: pagination.Page{
						Size:   10,
						Cursor: nil,
					},
				}
			},
			want: func(_ *testing.T, _ []entities.Transaction, _ *Repository) w {
				return w{
					entries: []vos.AccountEntry{},
					cursor:  nil,
				}
			},
		},
		{
			name: "return all entries",
			seedRepo: func(t *testing.T, ctx context.Context, r *Repository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), amount)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx := createTransaction(t, ctx, r, e1, e2)

				return []entities.Transaction{tx}
			},
			setupRequest: func(t *testing.T, _ []entities.Transaction, _ *Repository) vos.AccountEntryRequest {
				account, err := vos.NewAnalyticAccount(account1)
				assert.NoError(t, err)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Page: pagination.Page{
						Size:   10,
						Cursor: nil,
					},
				}
			},
			want: func(_ *testing.T, txs []entities.Transaction, _ *Repository) w {
				entries := accountEntriesFromTransaction(t, txs[0], account1)

				return w{
					entries: entries,
					cursor:  nil,
				}
			},
		},
		{
			name: "return entries from all accounts",
			seedRepo: func(t *testing.T, ctx context.Context, r *Repository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), amount)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx := createTransaction(t, ctx, r, e1, e2)

				return []entities.Transaction{tx}
			},
			setupRequest: func(t *testing.T, _ []entities.Transaction, _ *Repository) vos.AccountEntryRequest {
				account, err := vos.NewAccount(synthAccount)
				assert.NoError(t, err)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Page: pagination.Page{
						Size:   10,
						Cursor: nil,
					},
				}
			},
			want: func(_ *testing.T, txs []entities.Transaction, _ *Repository) w {
				entries := accountEntriesFromTransactionWithSyntheticAccountFilter(t, txs, synthAccount, 2)

				return w{
					entries: entries,
					cursor:  nil,
				}
			},
		},
		{
			name: "return first page",
			seedRepo: func(t *testing.T, ctx context.Context, r *Repository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), amount)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx1 := createTransaction(t, ctx, r, e1, e2)

				e1 = createEntry(t, vos.DebitOperation, account1, vos.Version(2), amount)
				e2 = createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx2 := createTransaction(t, ctx, r, e1, e2)

				return []entities.Transaction{tx1, tx2}
			},
			setupRequest: func(t *testing.T, _ []entities.Transaction, _ *Repository) vos.AccountEntryRequest {
				account, err := vos.NewAnalyticAccount(account1)
				assert.NoError(t, err)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Page: pagination.Page{
						Size:   1,
						Cursor: nil,
					},
				}
			},
			want: func(t *testing.T, txs []entities.Transaction, r *Repository) w {
				entries := accountEntriesFromTransaction(t, txs[1], account1)
				cur := cursorFromTransaction(t, txs[0], account1, r)

				return w{
					entries: entries,
					cursor:  cur,
				}
			},
		},
		{
			name: "return second page",
			seedRepo: func(t *testing.T, ctx context.Context, r *Repository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), amount)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx1 := createTransaction(t, ctx, r, e1, e2)

				e1 = createEntry(t, vos.DebitOperation, account1, vos.Version(2), amount)
				e2 = createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx2 := createTransaction(t, ctx, r, e1, e2)

				return []entities.Transaction{tx1, tx2}
			},
			setupRequest: func(t *testing.T, txs []entities.Transaction, r *Repository) vos.AccountEntryRequest {
				account, err := vos.NewAnalyticAccount(account1)
				assert.NoError(t, err)

				cur := cursorFromTransaction(t, txs[0], account1, r)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Page: pagination.Page{
						Size:   1,
						Cursor: cur,
					},
				}
			},
			want: func(t *testing.T, txs []entities.Transaction, _ *Repository) w {
				entries := accountEntriesFromTransaction(t, txs[0], account1)

				return w{
					entries: entries,
					cursor:  nil,
				}
			},
		},
		{
			name: "return filtered by single company",
			seedRepo: func(t *testing.T, ctx context.Context, r *Repository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), amount)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx1, err := entities.NewTransaction(
					uuid.New(),
					uint32(1),
					"company_1",
					time.Now().Round(time.Microsecond),
					e1, e2,
				)
				assert.NoError(t, err)

				err = r.CreateTransaction(ctx, tx1)
				assert.NoError(t, err)

				e1 = createEntry(t, vos.DebitOperation, account1, vos.Version(2), amount)
				e2 = createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx2 := createTransaction(t, ctx, r, e1, e2)

				return []entities.Transaction{tx1, tx2}
			},
			setupRequest: func(t *testing.T, _ []entities.Transaction, _ *Repository) vos.AccountEntryRequest {
				account, err := vos.NewAnalyticAccount(account1)
				assert.NoError(t, err)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Filter:    vos.AccountEntryFilter{Companies: []string{"company_1"}},
					Page: pagination.Page{
						Size:   10,
						Cursor: nil,
					},
				}
			},
			want: func(t *testing.T, txs []entities.Transaction, _ *Repository) w {
				entries := accountEntriesFromTransaction(t, txs[0], account1)

				return w{
					entries: entries,
					cursor:  nil,
				}
			},
		},
		{
			name: "return filtered by multiple companies",
			seedRepo: func(t *testing.T, ctx context.Context, r *Repository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), amount)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx1, err := entities.NewTransaction(
					uuid.New(),
					uint32(1),
					"company_1",
					time.Now().Round(time.Microsecond),
					e1, e2,
				)
				assert.NoError(t, err)

				err = r.CreateTransaction(ctx, tx1)
				assert.NoError(t, err)

				e1 = createEntry(t, vos.DebitOperation, account1, vos.Version(2), amount)
				e2 = createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx2 := createTransaction(t, ctx, r, e1, e2)

				return []entities.Transaction{tx1, tx2}
			},
			setupRequest: func(t *testing.T, _ []entities.Transaction, _ *Repository) vos.AccountEntryRequest {
				account, err := vos.NewAnalyticAccount(account1)
				assert.NoError(t, err)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Filter:    vos.AccountEntryFilter{Companies: []string{"company_1", "abc"}},
					Page: pagination.Page{
						Size:   10,
						Cursor: nil,
					},
				}
			},
			want: func(t *testing.T, txs []entities.Transaction, _ *Repository) w {
				entries := accountEntriesFromTransaction(t, txs[1], account1)
				entries = append(entries, accountEntriesFromTransaction(t, txs[0], account1)...)

				return w{
					entries: entries,
					cursor:  nil,
				}
			},
		},
		{
			name: "return filtered by single event",
			seedRepo: func(t *testing.T, ctx context.Context, r *Repository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), amount)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx1, err := entities.NewTransaction(
					uuid.New(),
					uint32(2),
					"abc",
					time.Now().Round(time.Microsecond),
					e1, e2,
				)
				assert.NoError(t, err)

				err = r.CreateTransaction(ctx, tx1)
				assert.NoError(t, err)

				e1 = createEntry(t, vos.DebitOperation, account1, vos.Version(2), amount)
				e2 = createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx2 := createTransaction(t, ctx, r, e1, e2)

				return []entities.Transaction{tx1, tx2}
			},
			setupRequest: func(t *testing.T, _ []entities.Transaction, _ *Repository) vos.AccountEntryRequest {
				account, err := vos.NewAnalyticAccount(account1)
				assert.NoError(t, err)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Filter:    vos.AccountEntryFilter{Events: []int32{2}},
					Page: pagination.Page{
						Size:   10,
						Cursor: nil,
					},
				}
			},
			want: func(t *testing.T, txs []entities.Transaction, _ *Repository) w {
				entries := accountEntriesFromTransaction(t, txs[0], account1)

				return w{
					entries: entries,
					cursor:  nil,
				}
			},
		},
		{
			name: "return filtered by multiple events",
			seedRepo: func(t *testing.T, ctx context.Context, r *Repository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), amount)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx1, err := entities.NewTransaction(
					uuid.New(),
					uint32(2),
					"abc",
					time.Now().Round(time.Microsecond),
					e1, e2,
				)
				assert.NoError(t, err)

				err = r.CreateTransaction(ctx, tx1)
				assert.NoError(t, err)

				e1 = createEntry(t, vos.DebitOperation, account1, vos.Version(2), amount)
				e2 = createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx2 := createTransaction(t, ctx, r, e1, e2)

				return []entities.Transaction{tx1, tx2}
			},
			setupRequest: func(t *testing.T, _ []entities.Transaction, _ *Repository) vos.AccountEntryRequest {
				account, err := vos.NewAnalyticAccount(account1)
				assert.NoError(t, err)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Filter:    vos.AccountEntryFilter{Events: []int32{1, 2}},
					Page: pagination.Page{
						Size:   10,
						Cursor: nil,
					},
				}
			},
			want: func(t *testing.T, txs []entities.Transaction, _ *Repository) w {
				entries := accountEntriesFromTransaction(t, txs[1], account1)
				entries = append(entries, accountEntriesFromTransaction(t, txs[0], account1)...)

				return w{
					entries: entries,
					cursor:  nil,
				}
			},
		},
		{
			name: "return filtered by operation",
			seedRepo: func(t *testing.T, ctx context.Context, r *Repository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), amount)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx1 := createTransaction(t, ctx, r, e1, e2)

				e1 = createEntry(t, vos.CreditOperation, account1, vos.Version(2), amount)
				e2 = createEntry(t, vos.DebitOperation, account2, vos.IgnoreAccountVersion, amount)

				tx2 := createTransaction(t, ctx, r, e1, e2)

				return []entities.Transaction{tx1, tx2}
			},
			setupRequest: func(t *testing.T, _ []entities.Transaction, _ *Repository) vos.AccountEntryRequest {
				account, err := vos.NewAnalyticAccount(account1)
				assert.NoError(t, err)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Filter:    vos.AccountEntryFilter{Operation: vos.DebitOperation},
					Page: pagination.Page{
						Size:   10,
						Cursor: nil,
					},
				}
			},
			want: func(t *testing.T, txs []entities.Transaction, _ *Repository) w {
				entries := accountEntriesFromTransaction(t, txs[0], account1)

				return w{
					entries: entries,
					cursor:  nil,
				}
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newDB(t, tt.name)

			r := NewRepository(db, &instrumentators.LedgerInstrumentator{})
			ctx := context.Background()

			txs := tt.seedRepo(t, ctx, r)

			req := tt.setupRequest(t, txs, r)

			resp, cur, err := r.ListAccountEntries(ctx, req)
			want := tt.want(t, txs, r)
			got := w{entries: resp, cursor: cur}

			assert.NoError(t, err)
			assert.Equal(t, len(want.entries), len(got.entries))
			for i := range got.entries {
				for j := range want.entries {
					if got.entries[i].ID == want.entries[j].ID {
						assert.WithinDuration(t, want.entries[j].CreatedAt, got.entries[i].CreatedAt, time.Minute)
						want.entries[j].CreatedAt = got.entries[i].CreatedAt
					}
				}
			}
			assert.Equal(t, want, got)
		})
	}
}

func TestLedgerRepository_ListAccountEntriesWithSyntheticAccount(t *testing.T) {
	t.Parallel()

	const (
		account1     = "asset.cash.cash"
		account2     = "asset.accounts_receivable.clients.1"
		account3     = "asset.accounts_receivable.clients.2"
		synthAccount = "asset.accounts_receivable.*"
	)

	testCases := []struct {
		name         string
		pageSize     int
		seedRepo     func(*testing.T, context.Context, *Repository) []entities.Transaction
		setupRequest func(*testing.T, []entities.Transaction, int) vos.AccountEntryRequest
		want         func(*testing.T, []entities.Transaction, int) []vos.AccountEntry
	}{
		{
			name:     "return last two entries",
			pageSize: 2,
			seedRepo: func(t *testing.T, ctx context.Context, r *Repository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), 100)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, 100)

				tx1 := createTransaction(t, ctx, r, e1, e2)

				e1 = createEntry(t, vos.DebitOperation, account1, vos.Version(2), 100)
				e2 = createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, 30)
				e3 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, 10)
				e4 := createEntry(t, vos.CreditOperation, account3, vos.IgnoreAccountVersion, 60)

				tx2 := createTransaction(t, ctx, r, e1, e2, e3, e4)

				return []entities.Transaction{tx1, tx2}
			},
			setupRequest: func(t *testing.T, txs []entities.Transaction, pageSize int) vos.AccountEntryRequest {
				account, err := vos.NewAccount(synthAccount)
				assert.NoError(t, err)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Page: pagination.Page{
						Size:   pageSize,
						Cursor: nil,
					},
				}
			},
			want: func(t *testing.T, txs []entities.Transaction, pageSize int) []vos.AccountEntry {
				return accountEntriesFromTransactionWithSyntheticAccountFilter(t, txs, synthAccount, pageSize)
			},
		},
		{
			name:     "multiple entries, return the last five",
			pageSize: 5,
			seedRepo: func(t *testing.T, ctx context.Context, r *Repository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account2, vos.IgnoreAccountVersion, 1)
				e2 := createEntry(t, vos.CreditOperation, account1, vos.Version(1), 1)

				tx1 := createTransaction(t, ctx, r, e1, e2)

				e1 = createEntry(t, vos.DebitOperation, account1, vos.Version(2), 100)
				e2 = createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, 30)
				e3 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, 10)
				e4 := createEntry(t, vos.CreditOperation, account3, vos.IgnoreAccountVersion, 60)

				tx2 := createTransaction(t, ctx, r, e1, e2, e3, e4)

				entries := make([]entities.Entry, 0, 20)
				for i := 0; i < 10; i++ {
					entries = append(entries,
						createEntry(t, vos.DebitOperation, account1, vos.NextAccountVersion, 100),
						createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, 100))
				}
				tx3 := createTransaction(t, ctx, r, entries...)

				return []entities.Transaction{tx1, tx2, tx3}
			},
			setupRequest: func(t *testing.T, txs []entities.Transaction, pageSize int) vos.AccountEntryRequest {
				account, err := vos.NewAccount(synthAccount)
				assert.NoError(t, err)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Page: pagination.Page{
						Size:   pageSize,
						Cursor: nil,
					},
				}
			},
			want: func(t *testing.T, txs []entities.Transaction, pageSize int) []vos.AccountEntry {
				return accountEntriesFromTransactionWithSyntheticAccountFilter(t, txs, synthAccount, pageSize)
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			db := newDB(t, tt.name)
			r := NewRepository(db, &instrumentators.LedgerInstrumentator{})
			ctx := context.Background()
			txs := tt.seedRepo(t, ctx, r)
			req := tt.setupRequest(t, txs, tt.pageSize)
			want := tt.want(t, txs, tt.pageSize)

			// execute
			var got []vos.AccountEntry
			for {
				resp, cur, err := r.ListAccountEntries(ctx, req)
				assert.NoError(t, err)
				if cur == nil {
					got = resp
					break
				}
				req.Page.Cursor = cur
			}

			// assert
			assert.Equal(t, len(want), len(got), "assert len")
			for i := range got {
				for j := range want {
					if got[i].ID == want[j].ID {
						assert.WithinDuration(t, want[j].CreatedAt, got[i].CreatedAt, time.Minute)
						want[j].CreatedAt = got[i].CreatedAt
					}
				}
			}
			assert.Equal(t, want, got)
		})
	}
}

func accountEntriesFromTransaction(t *testing.T, tx entities.Transaction, account string) []vos.AccountEntry {
	t.Helper()

	act := make([]vos.AccountEntry, 0, len(tx.Entries))
	for _, et := range tx.Entries {
		if et.Account.Value() != account {
			continue
		}

		var mt map[string]interface{}
		err := json.Unmarshal(et.Metadata, &mt)
		assert.NoError(t, err)

		act = append(act, vos.AccountEntry{
			ID:             et.ID,
			Account:        et.Account.Value(),
			Version:        et.Version,
			Operation:      et.Operation,
			Amount:         et.Amount,
			Event:          int(tx.Event),
			CreatedAt:      time.Now(),
			CompetenceDate: tx.CompetenceDate.Round(time.Microsecond),
			Metadata:       mt,
		})
	}

	return act
}

func accountEntriesFromTransactionWithSyntheticAccountFilter(t *testing.T, txs []entities.Transaction, syntheticAcc string, pageSize int) []vos.AccountEntry {
	t.Helper()

	regex := regexp.MustCompile(strings.ReplaceAll(syntheticAcc, "*", "[^.]+"))

	// order by competence date desc
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].CompetenceDate.After(txs[j].CompetenceDate)
	})

	var totalRows int
	for i := range txs {
		var entries []entities.Entry
		for _, e := range txs[i].Entries {
			if regex.MatchString(e.Account.Value()) {
				entries = append(entries, e)
			}
		}

		// order by id desc
		sort.Slice(entries, func(j, k int) bool {
			return entries[j].ID.String() > entries[k].ID.String()
		})

		txs[i].Entries = entries
		totalRows += len(entries)
	}

	act := make([]vos.AccountEntry, 0)
	var rowCount, startAtPosition int
	if totalRows > pageSize {
		totalPages := totalRows / pageSize
		if totalRows > totalPages*pageSize {
			startAtPosition = totalPages * pageSize
		} else {
			startAtPosition = totalPages * (pageSize - 1)
		}
	}

	for _, tx := range txs {
		for _, et := range tx.Entries {
			if rowCount < startAtPosition {
				rowCount++
				continue
			}
			var mt map[string]interface{}
			err := json.Unmarshal(et.Metadata, &mt)
			assert.NoError(t, err)

			act = append(act, vos.AccountEntry{
				ID:             et.ID,
				Account:        et.Account.Value(),
				Version:        et.Version,
				Operation:      et.Operation,
				Amount:         et.Amount,
				Event:          int(tx.Event),
				CreatedAt:      time.Now(),
				CompetenceDate: tx.CompetenceDate.Round(time.Microsecond),
				Metadata:       mt,
			})
		}
	}

	return act
}

func cursorFromTransaction(t *testing.T, tx entities.Transaction, account string, repository *Repository) pagination.Cursor {
	t.Helper()

	var et entities.Entry
	for _, entry := range tx.Entries {
		if entry.Account.Value() == account {
			et = entry
			break
		}
	}
	assert.NotEmpty(t, et)

	q := `select created_at from entry where id = $1`

	var createdAt time.Time
	err := repository.db.QueryRow(context.Background(), q, et.ID).Scan(&createdAt)
	assert.NoError(t, err)

	cur, err := pagination.NewCursor(listAccountEntriesCursor{
		ID:             et.ID.String(),
		CompetenceDate: tx.CompetenceDate.Round(time.Microsecond),
		CreatedAt:      createdAt,
		Version:        et.Version.AsInt64(),
	})
	assert.NoError(t, err)

	return cur
}
