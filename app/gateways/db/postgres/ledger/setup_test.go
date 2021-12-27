package ledger

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stone-co/the-amazing-ledger/app/tests/pgtesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	_, teardown, err := pgtesting.StartDockerContainer(pgtesting.DockerContainerConfig{
		DBName:  "ledger_test_database",
		Version: "13-alpine",
	})
	if err != nil {
		return 1
	}

	defer teardown()

	return m.Run()
}

func newDB(t *testing.T, name string) *pgxpool.Pool {
	pool := pgtesting.NewDB(t, name)

	_, err := pool.Exec(context.Background(), "insert into event (id, name) values (1, 'event_1'), (2, 'event_2');")
	require.NoError(t, err)

	return pool
}

func createEntry(t *testing.T, op vos.OperationType, account string, version vos.Version, amount int) entities.Entry {
	t.Helper()

	entry, err := entities.NewEntry(
		uuid.New(),
		op,
		account,
		version,
		amount,
		json.RawMessage(`{}`),
	)
	assert.NoError(t, err)

	return entry
}

func createTransaction(t *testing.T, ctx context.Context, r *Repository, entries ...entities.Entry) entities.Transaction {
	t.Helper()

	tx, err := entities.NewTransaction(
		uuid.New(),
		uint32(1),
		"abc",
		time.Now().Round(time.Microsecond),
		entries...,
	)
	assert.NoError(t, err)

	err = r.CreateTransaction(ctx, tx)
	assert.NoError(t, err)

	return tx
}
