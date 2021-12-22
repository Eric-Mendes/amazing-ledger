package usecases

import (
	"context"
	"fmt"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetAccountBalance(ctx context.Context, input domain.GetAccountBalanceInput) (vos.AccountBalance, error) {
	if !input.StartDate.IsZero() || !input.EndDate.IsZero() {
		return l.getBoundedAccountBalance(ctx, input)
	}

	return l.getRecentAccountBalance(ctx, input)
}

func (l *LedgerUseCase) getBoundedAccountBalance(ctx context.Context, input domain.GetAccountBalanceInput) (vos.AccountBalance, error) {
	balance, err := l.repository.GetBoundedAccountBalance(ctx, input.Account, input.StartDate, input.EndDate)
	if err != nil {
		return vos.AccountBalance{}, fmt.Errorf("get bounded account balance: %w", err)
	}

	return balance, nil
}

func (l *LedgerUseCase) getRecentAccountBalance(ctx context.Context, input domain.GetAccountBalanceInput) (vos.AccountBalance, error) {
	var (
		accountBalance vos.AccountBalance
		err            error
	)

	switch input.Account.Type() {
	case vos.Analytic:
		accountBalance, err = l.repository.GetAnalyticAccountBalance(ctx, input.Account)
	case vos.Synthetic:
		accountBalance, err = l.repository.GetSyntheticAccountBalance(ctx, input.Account)
	default:
		err = app.ErrInvalidAccountType
	}

	if err != nil {
		return vos.AccountBalance{}, fmt.Errorf("get account balance: %w", err)
	}

	return accountBalance, nil
}
