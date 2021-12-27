package ledger

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stone-co/the-amazing-ledger/app/domain/instrumentators"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestLedgerRepository_GetSyntheticReportSuccess(t *testing.T) {
	t.Parallel()

	const base = "liability.assets"

	creditAccount, err := vos.NewAccount(base + ".credit")
	require.NoError(t, err)

	debitAccount, err := vos.NewAccount(base + ".debit")
	require.NoError(t, err)

	testCases := []struct {
		name      string
		query     string
		level     int
		startDate time.Time
		endDate   time.Time
		seed      func(*testing.T, context.Context, *Repository)
		want      *vos.SyntheticReport
	}{
		{
			name:      "analytic account",
			query:     debitAccount.Value(),
			level:     3,
			startDate: time.Now().UTC().Add(-time.Second),
			endDate:   time.Now().UTC().Add(time.Hour),
			seed: func(t *testing.T, ctx context.Context, r *Repository) {
				e1 := createEntry(t, vos.DebitOperation, debitAccount.Value(), vos.IgnoreAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, creditAccount.Value(), vos.IgnoreAccountVersion, 100)

				createTransaction(t, ctx, r, e1, e2)
			},
			want: &vos.SyntheticReport{
				TotalCredit: 0,
				TotalDebit:  100,
				Results: []vos.AccountResult{
					{
						Account: debitAccount,
						Credit:  0,
						Debit:   100,
					},
				},
			},
		},
		{
			name:      "synthetic account",
			query:     base + ".*",
			level:     3,
			startDate: time.Now().UTC().Add(-time.Second),
			endDate:   time.Now().UTC().Add(time.Hour),
			seed: func(t *testing.T, ctx context.Context, r *Repository) {
				e1 := createEntry(t, vos.DebitOperation, debitAccount.Value(), vos.IgnoreAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, creditAccount.Value(), vos.IgnoreAccountVersion, 100)

				createTransaction(t, ctx, r, e1, e2)
			},
			want: &vos.SyntheticReport{
				TotalCredit: 100,
				TotalDebit:  100,
				Results: []vos.AccountResult{
					{
						Account: creditAccount,
						Credit:  100,
						Debit:   0,
					},
					{
						Account: debitAccount,
						Credit:  0,
						Debit:   100,
					},
				},
			},
		},
		{
			name:      "no data in database",
			query:     base + ".*",
			level:     3,
			startDate: time.Now().UTC(),
			endDate:   time.Now().UTC().Add(time.Hour * 1),
			seed:      func(_ *testing.T, _ context.Context, _ *Repository) {},
			want: &vos.SyntheticReport{
				TotalCredit: 0,
				TotalDebit:  0,
				Results:     nil,
			},
		},
		{
			name:      "no data for given account",
			query:     base + ".omni",
			level:     3,
			startDate: time.Now().UTC().Add(-time.Second),
			endDate:   time.Now().UTC().Add(time.Hour),
			seed: func(t *testing.T, ctx context.Context, r *Repository) {
				e1 := createEntry(t, vos.DebitOperation, debitAccount.Value(), vos.IgnoreAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, creditAccount.Value(), vos.IgnoreAccountVersion, 100)

				createTransaction(t, ctx, r, e1, e2)
			},
			want: &vos.SyntheticReport{
				TotalCredit: 0,
				TotalDebit:  0,
				Results:     nil,
			},
		},
		{
			name:      "no data in interval",
			query:     base + ".*",
			level:     3,
			startDate: time.Now().UTC().Add(-time.Hour),
			endDate:   time.Now().UTC().Add(-time.Second),
			seed: func(t *testing.T, ctx context.Context, r *Repository) {
				e1 := createEntry(t, vos.DebitOperation, debitAccount.Value(), vos.IgnoreAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, creditAccount.Value(), vos.IgnoreAccountVersion, 100)

				createTransaction(t, ctx, r, e1, e2)
			},
			want: &vos.SyntheticReport{
				TotalCredit: 0,
				TotalDebit:  0,
				Results:     nil,
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newDB(t, t.Name())
			ctx := context.Background()

			r := NewRepository(db, &instrumentators.LedgerInstrumentator{})
			tt.seed(t, ctx, r)

			query, err := vos.NewAccount(tt.query)
			assert.NoError(t, err)

			got, err := r.GetSyntheticReport(ctx, query, tt.level, tt.startDate, tt.endDate)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
