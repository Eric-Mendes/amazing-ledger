package ledger

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/instrumentation/newrelic"
)

const queryAggregatedBalanceQuery = `
select get_synthetic_account_balance($1);
`

func (r Repository) GetSyntheticAccountBalance(ctx context.Context, account vos.Account) (vos.AccountBalance, error) {
	const operation = "Repository.GetSyntheticAccountBalance"

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, queryAggregatedBalanceQuery).End()

	var balance int

	err := r.db.QueryRow(ctx, queryAggregatedBalanceQuery, account.Value()).Scan(&balance)
	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return vos.AccountBalance{}, fmt.Errorf("failed to query aggregated balance: %w", err)
		}

		if pgErr.Code == pgerrcode.NoDataFound {
			return vos.AccountBalance{}, app.ErrAccountNotFound
		}

		return vos.AccountBalance{}, fmt.Errorf("failed to query aggregated balance: %w", pgErr)
	}

	return vos.NewSyntheticAccountBalance(account, balance), nil
}
