package ledger

import (
	"context"
	"fmt"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/instrumentation/newrelic"
	pag "github.com/stone-co/the-amazing-ledger/app/pagination"
)

const (
	_accountEntriesQueryPrefix = `
select
	id,
	account,
	version,
	operation,
	amount,
	event,
	created_at,
	competence_date,
	metadata
from
	entry
where
	account %s $1
	and competence_date >= $2
	and competence_date < $3
`

	_accountEntriesCompanyFilter = `
	and company = $%d
`

	_accountEntriesCompaniesFilter = `
	and company = any($%d)
`

	_accountEntriesEventFilter = `
	and event = $%d
`

	_accountEntriesEventsFilter = `
	and event = any($%d)
`

	_accountEntriesOperationFilter = `
	and operation = $%d
`

	_accountEntriesQueryPaginationAnalytic = `
	and (competence_date, version) <= ($%d, $%d)
`
	_accountEntriesQueryPaginationSynthetic = `
	and ((competence_date <= $%d and created_at = $%d and id <= $%d)
	or (competence_date <= $%d and created_at < $%d))
`

	_accountEntriesQuerySuffixAnalytic = `
order by
	competence_date desc,
	version desc
limit $4;
`

	_accountEntriesQuerySuffixSynthetic = `
order by
	competence_date desc,
	created_at desc,
	id desc
limit $4;
`
)

type listAccountEntriesCursor struct {
	ID             string    `json:"id"`
	CompetenceDate time.Time `json:"competence_date"`
	CreatedAt      time.Time `json:"created_at"`
	Version        int64     `json:"version"`
}

func (r Repository) ListAccountEntries(ctx context.Context, req vos.AccountEntryRequest) ([]vos.AccountEntry, pag.Cursor, error) {
	const op = "Repository.ListAccountEntries"

	query, args, err := generateListAccountEntriesQuery(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate %s query: %w", op, err)
	}

	defer newrelic.NewDatastoreSegment(ctx, collection, op, query).End()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()

	entries := make([]vos.AccountEntry, 0)

	for rows.Next() {
		var entry vos.AccountEntry

		if err = rows.Scan(
			&entry.ID,
			&entry.Account,
			&entry.Version,
			&entry.Operation,
			&entry.Amount,
			&entry.Event,
			&entry.CreatedAt,
			&entry.CompetenceDate,
			&entry.Metadata,
		); err != nil {
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}

		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("%s rows have error: %w", op, err)
	}

	if len(entries) <= req.Page.Size {
		return entries, nil, nil
	}

	lastEntry := entries[len(entries)-1]
	entries = entries[:len(entries)-1]

	cursor, err := pag.NewCursor(listAccountEntriesCursor{
		ID:             lastEntry.ID.String(),
		CompetenceDate: lastEntry.CompetenceDate,
		CreatedAt:      lastEntry.CreatedAt,
		Version:        lastEntry.Version.AsInt64(),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate next page token: %w", err)
	}

	return entries, cursor, nil
}

func generateListAccountEntriesQuery(req vos.AccountEntryRequest) (string, []interface{}, error) {
	var (
		totalArgs = 4
		args      = []interface{}{req.Account.Value(), req.StartDate, req.EndDate, req.Page.Size + 1}
	)

	operator := "="
	if req.Account.Type() == vos.Synthetic {
		operator = "~"
	}

	query := fmt.Sprintf(_accountEntriesQueryPrefix, operator)

	switch len(req.Filter.Companies) {
	case 0:
		break
	case 1:
		query += fmt.Sprintf(_accountEntriesCompanyFilter, totalArgs+1)
		args = append(args, req.Filter.Companies[0])
		totalArgs += 1
	default:
		query += fmt.Sprintf(_accountEntriesCompaniesFilter, totalArgs+1)
		args = append(args, req.Filter.Companies)
		totalArgs += 1
	}

	switch len(req.Filter.Events) {
	case 0:
		break
	case 1:
		query += fmt.Sprintf(_accountEntriesEventFilter, totalArgs+1)
		args = append(args, req.Filter.Events[0])
		totalArgs += 1
	default:
		query += fmt.Sprintf(_accountEntriesEventsFilter, totalArgs+1)
		args = append(args, req.Filter.Events)
		totalArgs += 1
	}

	if req.Filter.Operation != vos.InvalidOperation {
		query += fmt.Sprintf(_accountEntriesOperationFilter, totalArgs+1)
		args = append(args, req.Filter.Operation)
		totalArgs += 1
	}

	if req.Page.Cursor != nil {
		var cursor listAccountEntriesCursor
		err := req.Page.Extract(&cursor)
		if err != nil {
			return "", nil, err
		}

		if req.Account.Type() == vos.Analytic {
			query += fmt.Sprintf(_accountEntriesQueryPaginationAnalytic, totalArgs+1, totalArgs+2)
			args = append(args, cursor.CompetenceDate, cursor.Version)
		}

		if req.Account.Type() == vos.Synthetic {
			query += fmt.Sprintf(_accountEntriesQueryPaginationSynthetic, totalArgs+1, totalArgs+2, totalArgs+3, totalArgs+1, totalArgs+2)
			args = append(args, cursor.CompetenceDate, cursor.CreatedAt, cursor.ID)
		}

	}

	if req.Account.Type() == vos.Analytic {
		query += _accountEntriesQuerySuffixAnalytic
	}

	if req.Account.Type() == vos.Synthetic {
		query += _accountEntriesQuerySuffixSynthetic
	}

	return query, args, nil
}
