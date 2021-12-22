package rpc

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger/v1beta"
)

func (a *API) GetAccountBalance(ctx context.Context, request *proto.GetAccountBalanceRequest) (*proto.GetAccountBalanceResponse, error) {
	accountName, err := vos.NewAccount(request.Account)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("can't create account name")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	start := time.Time{}
	if request.StartDate != nil && request.StartDate.IsValid() {
		start = request.StartDate.AsTime()
	}

	end := time.Time{}
	if request.EndDate != nil && request.EndDate.IsValid() {
		end = request.EndDate.AsTime()
	}

	if !start.IsZero() && !end.IsZero() && end.Before(start) {
		return nil, status.Error(codes.InvalidArgument, "end date should be a timestamp set after start date")
	}

	input := domain.GetAccountBalanceInput{
		Account:   accountName,
		StartDate: start,
		EndDate:   end,
	}

	accountBalance, err := a.UseCase.GetAccountBalance(ctx, input)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to get account balance")
		if errors.Is(err, app.ErrAccountNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &proto.GetAccountBalanceResponse{
		Account:        accountBalance.Account.Value(),
		CurrentVersion: accountBalance.CurrentVersion.AsInt64(),
		Balance:        int64(accountBalance.Balance),
	}, nil
}
