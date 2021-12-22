package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests/mocks"
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger/v1beta"
)

func TestAPI_GetAccountBalance_Analytic_Success(t *testing.T) {
	t.Run("should get account balance successfully", func(t *testing.T) {
		account, err := vos.NewAccount(testdata.GenerateAccountPath())
		assert.NoError(t, err)

		accountBalance := vos.NewAnalyticAccountBalance(account, vos.Version(1), 200)
		mockedUsecase := &mocks.UseCaseMock{
			GetAccountBalanceFunc: func(ctx context.Context, input domain.GetAccountBalanceInput) (vos.AccountBalance, error) {
				accountBalance.Account = input.Account

				return accountBalance, nil
			},
		}
		api := NewAPI(mockedUsecase)

		request := &proto.GetAccountBalanceRequest{
			Account: account.Value(),
		}

		got, err := api.GetAccountBalance(context.Background(), request)
		assert.NoError(t, err)

		assert.Equal(t, &proto.GetAccountBalanceResponse{
			Account:        request.Account,
			CurrentVersion: accountBalance.CurrentVersion.AsInt64(),
			Balance:        200,
		}, got)
	})
}

func TestAPI_GetAccountBalance_Synthetic_Success(t *testing.T) {
	t.Run("should get aggregated balance successfully", func(t *testing.T) {
		account, err := vos.NewAccount("liability.stone.clients.*")
		assert.NoError(t, err)

		balance := vos.NewSyntheticAccountBalance(account, 100)
		mockedUsecase := &mocks.UseCaseMock{
			GetAccountBalanceFunc: func(ctx context.Context, input domain.GetAccountBalanceInput) (vos.AccountBalance, error) {
				return balance, nil
			},
		}
		api := NewAPI(mockedUsecase)

		request := &proto.GetAccountBalanceRequest{
			Account: "liability.stone.clients.*",
		}

		got, err := api.GetAccountBalance(context.Background(), request)
		assert.NoError(t, err)

		assert.Equal(t, &proto.GetAccountBalanceResponse{
			Account:        account.Value(),
			CurrentVersion: -1,
			Balance:        100,
		}, got)
	})
}

func TestAPI_GetAccountBalance_Bounds(t *testing.T) {
	t.Parallel()

	account, err := vos.NewAccount("liability.stone.clients.*")
	require.NoError(t, err)

	balance := vos.NewSyntheticAccountBalance(account, 100)
	mockedUsecase := &mocks.UseCaseMock{
		GetAccountBalanceFunc: func(ctx context.Context, input domain.GetAccountBalanceInput) (vos.AccountBalance, error) {
			return balance, nil
		},
	}
	api := NewAPI(mockedUsecase)

	tests := []struct {
		name    string
		request *proto.GetAccountBalanceRequest
	}{
		{
			name: "both dates set",
			request: &proto.GetAccountBalanceRequest{
				Account:   "liability.stone.clients.*",
				StartDate: timestamppb.Now(),
				EndDate:   timestamppb.New(time.Now().Add(1 * time.Second)),
			},
		},
		{
			name: "only start set",
			request: &proto.GetAccountBalanceRequest{
				Account:   "liability.stone.clients.*",
				StartDate: timestamppb.Now(),
			},
		},
		{
			name: "only end set",
			request: &proto.GetAccountBalanceRequest{
				Account: "liability.stone.clients.*",
				EndDate: timestamppb.Now(),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := api.GetAccountBalance(context.Background(), tt.request)
			assert.NoError(t, err)

			assert.Equal(t, &proto.GetAccountBalanceResponse{
				Account:        account.Value(),
				CurrentVersion: -1,
				Balance:        100,
			}, got)
		})
	}
}

func TestAPI_GetAccountBalance_InvalidRequest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		useCaseSetup    *mocks.UseCaseMock
		request         *proto.GetAccountBalanceRequest
		expectedCode    codes.Code
		expectedMessage string
	}{
		{
			name:         "should return an error if account name is invalid",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.GetAccountBalanceRequest{
				Account: "liability.clients.abc-123.*",
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: app.ErrInvalidAccountComponentCharacters.Error(),
		},
		{
			name:         "should return an error if dates are invalid",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.GetAccountBalanceRequest{
				Account:   testdata.GenerateAccount(),
				StartDate: timestamppb.Now(),
				EndDate:   timestamppb.New(time.Now().Add(-1 * time.Second)),
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "end date should be a timestamp set after start date",
		},
		{
			name: "should return an error if account does not exist",
			useCaseSetup: &mocks.UseCaseMock{
				GetAccountBalanceFunc: func(ctx context.Context, input domain.GetAccountBalanceInput) (vos.AccountBalance, error) {
					return vos.AccountBalance{}, app.ErrAccountNotFound
				},
			},
			request: &proto.GetAccountBalanceRequest{
				Account: testdata.GenerateAccountPath(),
			},
			expectedCode:    codes.NotFound,
			expectedMessage: "account not found",
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			api := NewAPI(tt.useCaseSetup)

			_, err := api.GetAccountBalance(context.Background(), tt.request)
			respStatus, ok := status.FromError(err)

			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, respStatus.Code())
			assert.Equal(t, tt.expectedMessage, respStatus.Message())
		})
	}
}
