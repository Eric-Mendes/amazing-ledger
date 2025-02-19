package ledger

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/instrumentation/newrelic"
)

const (
	numArgs           = 10
	numDefaultQueries = 5
)

const createTransactionQuery = `
insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company, metadata)
values %s;`

func (r Repository) CreateTransaction(ctx context.Context, transaction entities.Transaction) error {
	const operation = "Repository.CreateTransaction"

	query := r.qb.Build(len(transaction.Entries))
	args := make([]interface{}, 0)

	for _, entry := range transaction.Entries {
		args = append(
			args,
			entry.ID,
			transaction.ID,
			transaction.Event,
			entry.Operation,
			entry.Version,
			entry.Amount,
			transaction.CompetenceDate,
			entry.Account.Value(),
			transaction.Company,
			entry.Metadata,
		)
	}

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, query).End()

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); !ok {
			return err
		}

		switch pgErr.Code {
		case pgerrcode.RaiseException:
			return app.ErrInvalidVersion
		case pgerrcode.UniqueViolation:
			return app.ErrIdempotencyKeyViolation
		default:
			return err
		}
	}

	return nil
}
