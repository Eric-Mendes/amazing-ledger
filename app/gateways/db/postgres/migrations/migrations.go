package migrations

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var migrations embed.FS

func GetMigrationHandler(dbUrl string) (*migrate.Migrate, error) {
	source, err := iofs.New(migrations, ".")
	if err != nil {
		return nil, fmt.Errorf("init iofs: %w", err)
	}

	handler, err := migrate.NewWithSourceInstance("iofs", source, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("creating migration source: %w", err)
	}

	return handler, nil
}

func RunMigrations(connString string) error {
	handler, err := GetMigrationHandler(connString)
	if err != nil {
		return fmt.Errorf("get migration handler: %w", err)
	}

	defer handler.Close()

	if err = handler.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("running migrations: %w", err)
	}

	return nil
}
