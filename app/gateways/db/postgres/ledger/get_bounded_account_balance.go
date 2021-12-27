package ledger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/instrumentation/newrelic"
)

const _boundedBalanceQuery = `
select coalesce(sum(amount) filter (where operation = 1), 0) -
	   coalesce(sum(amount) filter (where operation = 2), 0)
  from entry
 where account %s $1
`

const _boundedBalanceQueryStartFilter = `
   and competence_date >= $%d
`

const _boundedBalanceQueryEndFilter = `
   and competence_date <  $%d
`

func (r Repository) GetBoundedAccountBalance(ctx context.Context, acc vos.Account, start, end time.Time) (vos.AccountBalance, error) {
	const operation = "Repository.GetBoundedAccountBalance"

	query, args := buildBoundedBalanceQuery(acc, start, end)

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, query).End()

	var balance int

	err := r.db.QueryRow(ctx, query, args...).Scan(&balance)
	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return vos.AccountBalance{}, fmt.Errorf("get account balance: %w", err)
		}

		if pgErr.Code == pgerrcode.NoDataFound {
			return vos.AccountBalance{}, app.ErrAccountNotFound
		}

		return vos.AccountBalance{}, fmt.Errorf("get account balance: %w", pgErr)
	}

	return vos.NewSyntheticAccountBalance(acc, balance), nil
}

func buildBoundedBalanceQuery(account vos.Account, start, end time.Time) (string, []interface{}) {
	operator := "="
	if account.Type() == vos.Synthetic {
		operator = "~"
	}

	args := make([]interface{}, 0, 3)
	args = append(args, account.Value())
	total := 2

	query := fmt.Sprintf(_boundedBalanceQuery, operator)

	if !start.IsZero() {
		query += fmt.Sprintf(_boundedBalanceQueryStartFilter, total)
		total += 1
		args = append(args, start)
	}

	if !end.IsZero() {
		query += fmt.Sprintf(_boundedBalanceQueryEndFilter, total)
		args = append(args, end)
	}

	return query, args
}
