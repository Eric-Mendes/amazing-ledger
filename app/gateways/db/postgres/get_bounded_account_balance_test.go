package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/instrumentators"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests"
)

func TestLedgerRepository_QueryBoundedBalance(t *testing.T) {
	ctx := context.Background()
	r := NewLedgerRepository(pgDocker.DB, &instrumentators.LedgerInstrumentator{})

	acc1, err := vos.NewAccount("liability.agg.agg1")
	assert.NoError(t, err)

	acc2, err := vos.NewAccount("liability.agg.agg2")
	assert.NoError(t, err)

	acc3, err := vos.NewAccount("liability.abc.agg3")
	assert.NoError(t, err)

	agg, err := vos.NewAccount("liability.agg.*")
	assert.NoError(t, err)

	e1 := createEntry(t, vos.DebitOperation, acc1.Value(), vos.NextAccountVersion, 100)
	e2 := createEntry(t, vos.CreditOperation, acc2.Value(), vos.IgnoreAccountVersion, 100)
	createTransactionWithDate(t, ctx, r, time.Now().Add(-48*time.Hour), e1, e2)

	e1 = createEntry(t, vos.CreditOperation, acc1.Value(), vos.NextAccountVersion, 200)
	e2 = createEntry(t, vos.CreditOperation, acc2.Value(), vos.NextAccountVersion, 100)
	e3 := createEntry(t, vos.DebitOperation, acc3.Value(), vos.IgnoreAccountVersion, 300)
	createTransactionWithDate(t, ctx, r, time.Now().Add(-24*time.Hour), e1, e2, e3)

	defer tests.TruncateTables(ctx, pgDocker.DB, "entry", "account_version", "account_balance")

	testCases := []struct {
		name    string
		account vos.Account
		start   time.Time
		end     time.Time
		wants   int
	}{
		{
			name:    "start date only -> last transaction only",
			account: acc1,
			start:   time.Now().Add(-25 * time.Hour),
			wants:   200,
		},
		{
			name:    "start and end date -> all transactions",
			account: acc1,
			start:   time.Now().Add(-49 * time.Hour),
			end:     time.Now().Add(-23 * time.Hour),
			wants:   100,
		},
		{
			name:    "end date only -> first transaction",
			account: acc1,
			end:     time.Now().Add(-25 * time.Hour),
			wants:   -100,
		},
		{
			name:    "start and end date -> synthetic account",
			account: agg,
			start:   time.Now().Add(-49 * time.Hour),
			end:     time.Now().Add(-23 * time.Hour),
			wants:   300,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			balance, err := r.GetBoundedAccountBalance(ctx, tt.account, tt.start, tt.end)
			assert.NoError(t, err)
			assert.Equal(t, tt.wants, balance.Balance)
		})
	}
}

func createTransactionWithDate(t *testing.T, ctx context.Context, r *LedgerRepository, dt time.Time, entries ...entities.Entry) entities.Transaction {
	t.Helper()

	tx, err := entities.NewTransaction(
		uuid.New(),
		uint32(1),
		"company",
		dt,
		entries...,
	)
	assert.NoError(t, err)

	err = r.CreateTransaction(ctx, tx)
	assert.NoError(t, err)

	return tx
}
