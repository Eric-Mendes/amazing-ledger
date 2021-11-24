package testenv

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/stone-co/the-amazing-ledger/app/gateways/db/postgres"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger/v1beta"
)

var (
	DB               *pgxpool.Pool
	RPCClient        proto.LedgerAPIClient
	LedgerRepository *postgres.LedgerRepository
	GatewayServer    string
)
