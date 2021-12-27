package pgtesting

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

// shared pool used for creating and connecting to new databases. It acts like a "mother-connection"
// it is ok to be kept as a package variable as this is intended to be used in tests and
// different test packages run in different PIDs, so they won't share this globally.
var concurrentConn *pgxpool.Pool

var nonAlphaRegex = regexp.MustCompile(`[\W]`)

// NewDB creates a new database named as a sanitized dbName. It returns a connection pool to this database.
// It must be called after StartDockerContainer.
func NewDB(t *testing.T, dbName string) *pgxpool.Pool {
	if concurrentConn == nil {
		return nil
	}
	t.Helper()

	if dbName == "" {
		require.FailNow(t, "dbName cannot be an empty string")
	}

	dbName = nonAlphaRegex.ReplaceAllString(strings.ToLower(dbName), "_")

	pool := concurrentConn
	// Just ensuring we didn't panic last time or anything.
	_, _ = pool.Exec(context.Background(), fmt.Sprintf("drop database %s", dbName))

	_, err := pool.Exec(context.Background(), fmt.Sprintf("create database %s TEMPLATE %s_template", dbName, pool.Config().ConnConfig.Database))
	require.NoError(t, err)

	connString := strings.Replace(pool.Config().ConnString(), pool.Config().ConnConfig.Database, dbName, 1)
	newPool, err := pgxpool.Connect(context.Background(), connString)
	require.NoError(t, err)

	t.Cleanup(func() {
		newPool.Close()
		_, _ = pool.Exec(context.Background(), fmt.Sprintf("drop database %s", dbName))
	})

	return newPool
}
