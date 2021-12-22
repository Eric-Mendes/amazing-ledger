package domain

import (
	"context"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/pagination"
)

type Repository interface {
	CreateTransaction(context.Context, entities.Transaction) error
	GetBoundedAccountBalance(context.Context, vos.Account, time.Time, time.Time) (vos.AccountBalance, error)
	GetAnalyticAccountBalance(context.Context, vos.Account) (vos.AccountBalance, error)
	GetSyntheticAccountBalance(context.Context, vos.Account) (vos.AccountBalance, error)
	GetSyntheticReport(context.Context, vos.Account, int, time.Time, time.Time) (*vos.SyntheticReport, error)
	ListAccountEntries(context.Context, vos.AccountEntryRequest) ([]vos.AccountEntry, pagination.Cursor, error)
}
