package rpc

import (
	"github.com/stone-co/the-amazing-ledger/app/domain"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger/v1beta"
)

var _ proto.LedgerAPIServer = &API{}

type API struct {
	UseCase domain.UseCase
}

func NewAPI(useCase domain.UseCase) *API {
	return &API{
		UseCase: useCase,
	}
}
