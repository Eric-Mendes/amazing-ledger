package domain

import (
	"context"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

type UseCase interface {
	CreateTransaction(context.Context, entities.Transaction) error
	GetAccountBalance(context.Context, GetAccountBalanceInput) (vos.AccountBalance, error)
	GetSyntheticReport(context.Context, vos.Account, int, time.Time, time.Time) (*vos.SyntheticReport, error)
	ListAccountEntries(context.Context, vos.AccountEntryRequest) (vos.AccountEntryResponse, error)
}

type GetAccountBalanceInput struct {
	Account   vos.Account
	StartDate time.Time
	EndDate   time.Time
}
