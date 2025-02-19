package ledger

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/instrumentators"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestLedgerRepository_CreateTransactionSuccess(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                   string
		repoSeed               func(t *testing.T, ctx context.Context, r *Repository)
		entriesSetup           func(t *testing.T) []entities.Entry
		expectedEntryVersion   vos.Version
		expectedAccountVersion vos.Version
	}{
		{
			name:     "insert transaction successfully with no previous versions - auto version",
			repoSeed: func(t *testing.T, ctx context.Context, r *Repository) {},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion, 100)

				return []entities.Entry{e1, e2}
			},
			expectedEntryVersion:   vos.Version(1),
			expectedAccountVersion: vos.Version(1),
		},
		{
			name:     "insert transaction successfully with no previous versions - manual version",
			repoSeed: func(t *testing.T, ctx context.Context, r *Repository) {},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.Version(1), 100)
				e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion, 100)

				return []entities.Entry{e1, e2}
			},
			expectedEntryVersion:   vos.Version(1),
			expectedAccountVersion: vos.Version(1),
		},
		{
			name: "insert transaction successfully with existing versions - auto version",
			repoSeed: func(t *testing.T, ctx context.Context, r *Repository) {
				e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion, 100)

				createTransaction(t, ctx, r, e1, e2)
			},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion, 100)

				return []entities.Entry{e1, e2}
			},
			expectedEntryVersion:   vos.Version(2),
			expectedAccountVersion: vos.Version(2),
		},
		{
			name: "insert transaction successfully with existing versions - manual version",
			repoSeed: func(t *testing.T, ctx context.Context, r *Repository) {
				e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion, 100)

				createTransaction(t, ctx, r, e1, e2)
			},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.Version(2), 100)
				e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion, 100)

				return []entities.Entry{e1, e2}
			},
			expectedEntryVersion:   vos.Version(2),
			expectedAccountVersion: vos.Version(2),
		},
		{
			name: "insert transaction successfully with existing versions - unsorted and manual version",
			repoSeed: func(t *testing.T, ctx context.Context, r *Repository) {
				e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion, 100)

				createTransaction(t, ctx, r, e1, e2)
			},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.Version(2), 40)
				e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion, 100)
				e3 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.Version(4), 40)
				e4 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.Version(3), 20)

				return []entities.Entry{e1, e2, e3, e4}
			},
			expectedEntryVersion:   vos.Version(2),
			expectedAccountVersion: vos.Version(4),
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newDB(t, tt.name)

			ctx := context.Background()
			r := NewRepository(db, &instrumentators.LedgerInstrumentator{})

			tt.repoSeed(t, ctx, r)

			entries := tt.entriesSetup(t)
			e1, e2 := entries[0], entries[1]

			tx, err := entities.NewTransaction(uuid.New(), uint32(1), "abc", time.Now(), entries...)
			assert.NoError(t, err)

			err = r.CreateTransaction(ctx, tx)
			assert.NoError(t, err)

			assertMetadata(t, ctx, db, e1.ID, e1.Metadata)
			assertMetadata(t, ctx, db, e2.ID, e2.Metadata)

			assertEntryVersion(t, ctx, db, e1.ID, tt.expectedEntryVersion)
			assertEntryVersion(t, ctx, db, e2.ID, vos.IgnoreAccountVersion)

			assertAccountVersion(t, ctx, db, e1.Account, tt.expectedAccountVersion)
			assertAccountVersion(t, ctx, db, e2.Account, vos.Version(0))
		})
	}
}

func TestLedgerRepository_CreateTransactionFailure(t *testing.T) {
	t.Parallel()

	e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.NextAccountVersion, 100)
	e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion, 100)

	testCases := []struct {
		name                   string
		repoSeed               func(t *testing.T, ctx context.Context, r *Repository)
		entriesSetup           func(t *testing.T) []entities.Entry
		expectedErr            error
		expectedAccountVersion vos.Version
	}{
		{
			name: "same version",
			repoSeed: func(t *testing.T, ctx context.Context, r *Repository) {
				createTransaction(t, ctx, r, e1, e2)
			},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e3 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.Version(1), 100)
				e4 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion, 100)

				return []entities.Entry{e3, e4}
			},
			expectedErr:            app.ErrInvalidVersion,
			expectedAccountVersion: vos.Version(1),
		},
		{
			name: "lower version",
			repoSeed: func(t *testing.T, ctx context.Context, r *Repository) {
				createTransaction(t, ctx, r, e1, e2)

				e3 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.NextAccountVersion, 100)
				e4 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion, 100)

				createTransaction(t, ctx, r, e3, e4)
			},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e5 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.Version(1), 100)
				e6 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion, 100)

				return []entities.Entry{e5, e6}
			},
			expectedErr:            app.ErrInvalidVersion,
			expectedAccountVersion: vos.Version(2),
		},
		{
			name: "random high version",
			repoSeed: func(t *testing.T, ctx context.Context, r *Repository) {
				createTransaction(t, ctx, r, e1, e2)
			},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e3 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.Version(30), 100)
				e4 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion, 100)

				return []entities.Entry{e3, e4}
			},
			expectedErr:            app.ErrInvalidVersion,
			expectedAccountVersion: vos.Version(1),
		},
		{
			name: "reusing entry id",
			repoSeed: func(t *testing.T, ctx context.Context, r *Repository) {
				createTransaction(t, ctx, r, e1, e2)
			},
			entriesSetup: func(t *testing.T) []entities.Entry {
				return []entities.Entry{e1, e2}
			},
			expectedErr:            app.ErrIdempotencyKeyViolation,
			expectedAccountVersion: vos.Version(1),
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newDB(t, t.Name())

			ctx := context.Background()
			r := NewRepository(db, &instrumentators.LedgerInstrumentator{})

			tt.repoSeed(t, ctx, r)

			entries := tt.entriesSetup(t)
			e1, e2 := entries[0], entries[1]

			tx, err := entities.NewTransaction(uuid.New(), uint32(1), "abc", time.Now(), entries...)
			assert.NoError(t, err)

			err = r.CreateTransaction(ctx, tx)
			assert.ErrorIs(t, err, tt.expectedErr)

			assertAccountVersion(t, ctx, db, e1.Account, tt.expectedAccountVersion)
			assertAccountVersion(t, ctx, db, e2.Account, vos.Version(0))
		})
	}
}

func assertAccountVersion(t *testing.T, ctx context.Context, db *pgxpool.Pool, account vos.Account, want vos.Version) {
	t.Helper()

	const query = `select coalesce(version, 0) from account_version where account = $1;`

	var version vos.Version

	err := db.QueryRow(ctx, query, account.Value()).Scan(&version)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("unexpected error: %v", err)
	} else if errors.Is(err, pgx.ErrNoRows) {
		assert.Equal(t, vos.Version(0), version)
	}

	assert.Equal(t, want, version)
}

func assertEntryVersion(t *testing.T, ctx context.Context, db *pgxpool.Pool, id uuid.UUID, want vos.Version) {
	t.Helper()

	const query = `select version from entry where id = $1;`

	var version vos.Version

	err := db.QueryRow(ctx, query, id).Scan(&version)
	require.NoError(t, err)

	assert.Equal(t, want, version)
}

func assertMetadata(t *testing.T, ctx context.Context, db *pgxpool.Pool, id uuid.UUID, want json.RawMessage) {
	t.Helper()

	const query = `select metadata from entry where id = $1`

	var metadata json.RawMessage

	err := db.QueryRow(ctx, query, id).Scan(&metadata)
	require.NoError(t, err)

	assert.Equal(t, string(want), string(metadata))
}
