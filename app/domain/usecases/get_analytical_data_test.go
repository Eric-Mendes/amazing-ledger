package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestLedgerUseCase_GetAnalyticalData(t *testing.T) {
	res := []vos.Statement{}
	fn := func(st vos.Statement) error {
		res = append(res, st)
		return nil
	}

	t.Run("Account can be empty", func(t *testing.T) {
		query, err := vos.NewAccountQuery("liability.stone.clients")
		assert.Nil(t, err)

		entries := []vos.Statement{}
		useCase := newFakeGetAnalyticalData(entries, nil)

		res = []vos.Statement{}
		err = useCase.GetAnalyticalData(context.Background(), query, fn)
		assert.Nil(t, err)
		assert.Equal(t, entries, res)
	})

	t.Run("Account with 2 entries", func(t *testing.T) {
		query, err := vos.NewAccountQuery("liability.stone.clients")
		assert.Nil(t, err)

		accountName1, err := vos.NewAccountPath("liability.stone.clients.user_1")
		assert.NoError(t, err)
		accountName2, err := vos.NewAccountPath("liability.stone.clients.user_2")
		assert.NoError(t, err)

		entries := []vos.Statement{
			{
				Account:   accountName1.Name(),
				Operation: vos.CreditOperation,
				Amount:    15000,
			},
			{
				Account:   accountName2.Name(),
				Operation: vos.CreditOperation,
				Amount:    25000,
			},
		}

		useCase := newFakeGetAnalyticalData(entries, nil)

		res = []vos.Statement{}
		err = useCase.GetAnalyticalData(context.Background(), query, fn)
		assert.Nil(t, err)
		assert.Equal(t, entries, res)
	})
}
