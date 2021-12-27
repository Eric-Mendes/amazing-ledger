package ledger

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/instrumentators"
	"github.com/stone-co/the-amazing-ledger/app/gateways/db/querybuilder"
)

const collection = "entry"

var _ domain.Repository = &Repository{}

type Repository struct {
	db *pgxpool.Pool
	pb *instrumentators.LedgerInstrumentator
	qb querybuilder.QueryBuilder
}

func NewRepository(db *pgxpool.Pool, pb *instrumentators.LedgerInstrumentator) *Repository {
	qb := querybuilder.New(createTransactionQuery, numArgs)
	qb.Init(numDefaultQueries)

	return &Repository{
		db: db,
		pb: pb,
		qb: qb,
	}
}
