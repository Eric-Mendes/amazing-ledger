package rpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests"
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
	"github.com/stone-co/the-amazing-ledger/app/tests/testenv"
	"github.com/stone-co/the-amazing-ledger/app/tests/testseed"
	"github.com/stone-co/the-amazing-ledger/app/tests/testutils"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func TestE2E_RPC_GetAccountBalanceSuccess(t *testing.T) {
	t.Run("should get account balance successfully", func(t *testing.T) {
		e1 := testutils.CreateEntry(t, vos.DebitOperation, testdata.GenerateAccountPath(), vos.NextAccountVersion, 100)
		e2 := testutils.CreateEntry(t, vos.CreditOperation, testdata.GenerateAccountPath(), vos.NextAccountVersion, 100)

		_ = testseed.CreateTransaction(t, e1, e2)

		defer tests.TruncateTables(context.Background(), testenv.DB, "entry", "account_version")

		request := &proto.GetAccountBalanceRequest{
			AccountPath: e1.Account.Value(),
		}

		balance, err := testenv.RPCClient.GetAccountBalance(context.Background(), request)
		assert.NoError(t, err)
		assert.Equal(t, int64(-100), balance.Balance)

		request = &proto.GetAccountBalanceRequest{
			AccountPath: e2.Account.Value(),
		}

		balance, err = testenv.RPCClient.GetAccountBalance(context.Background(), request)
		assert.NoError(t, err)
		assert.Equal(t, int64(100), balance.Balance)
	})
}

func TestE2E_RPC_GetAccountBalanceFailure(t *testing.T) {
	testCases := []struct {
		name         string
		account      string
		expectedCode codes.Code
		expectedMsg  string
	}{
		{
			name:         "should return an error if account does not exist",
			account:      testdata.GenerateAccountPath(),
			expectedCode: codes.NotFound,
			expectedMsg:  "failed to get account balance: account not found",
		},
		{
			name:         "should return an error if account path is invalid",
			account:      "liability.asset",
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "account does not meet minimum or maximum supported sizes",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			request := &proto.GetAccountBalanceRequest{
				AccountPath: tt.account,
			}

			balance, err := testenv.RPCClient.GetAccountBalance(context.Background(), request)
			assert.Nil(t, balance)

			status, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, status.Code())
			assert.Equal(t, tt.expectedMsg, status.Message())
		})
	}
}