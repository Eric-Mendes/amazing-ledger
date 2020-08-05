package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts/usecase"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/db/postgres"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/accounts"
)

func main() {
	log := logrus.New()
	log.Infoln("Starting Command-Handler process...")

	cfg, err := configuration.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("unable to load app configuration: %s", err.Error()))
	}

	conn, err := postgres.ConnectPool(cfg.Postgres.DSN(), log)
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err.Error())
	}
	defer conn.Close()

	if err := postgres.RunMigrations(cfg.Postgres.URL()); err != nil {
		log.WithError(err).Fatal("Error running postgres migrations")
	}

	accountsRepository := postgres.NewAccountsRepository(conn, log)
	accountsUseCase := usecase.NewAccountUseCase(log, accountsRepository)
	accountsHandler := accounts.NewAccountsHandler(log, accountsUseCase)

	// Starting gateway http API
	api := http.NewApi(log, accountsHandler)
	api.Start("0.0.0.0", cfg.API.Port)
}
