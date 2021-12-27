package pgtesting

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog"
	"github.com/stone-co/the-amazing-ledger/app/gateways/db/postgres"
	"github.com/stone-co/the-amazing-ledger/app/gateways/db/postgres/migrations"
)

const (
	_postgresDefaultDB            = "postgres"
	_defaultPostgresDockerVersion = "13-alpine"
)

type DockerContainerConfig struct {
	// DBName is the name of the postgres database.
	DBName string

	// Version is the Postgres version ran in the container (default is 13).
	Version string

	// Expire container that takes more than `Expire` seconds running.
	// It is very important to use when you are debugging and killed the
	// process before the call the teardown or if a panic happens.
	Expire uint
}

var instances int32

type DockerizedPostgres struct {
	Pool *pgxpool.Pool
	Port string
}

func StartDockerContainer(cfg DockerContainerConfig) (_ *DockerizedPostgres, teardownFn func(), err error) {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf(`could not connect to docker: %w`, err)
	}

	if err = dockerPool.Client.Ping(); err != nil {
		return nil, nil, fmt.Errorf(`could not connect to docker: %w`, err)
	}

	if cfg.Version == "" {
		cfg.Version = _defaultPostgresDockerVersion
	}

	if cfg.DBName == "" {
		atomic.AddInt32(&instances, 1)
		cfg.DBName = fmt.Sprintf("db_%d_%d", time.Now().UnixNano(), atomic.LoadInt32(&instances))
	}

	dockerResource, err := getDockerPostgresResource(dockerPool, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf(`failed to initialize postgres docker resource: %w`, err)
	}

	dbPort := dockerResource.GetPort("5432/tcp")

	// Container started but postgres might still be starting;
	// keep trying to connect until connection is established.
	if err = dockerPool.Retry(pingPostgresFn(dbPort)); err != nil {
		return nil, nil, err
	}

	if cfg.Expire != 0 {
		_ = dockerResource.Expire(cfg.Expire)
	}

	defaultPGPool, err := postgres.ConnectPool(
		context.Background(),
		getPostgresConnString(dockerResource.GetPort("5432/tcp"), "postgres"),
		zerolog.Nop(),
	)
	if err != nil {
		return nil, nil, err
	}

	// Creates and connects to the database that will be used for tests.
	if err = createDB(cfg.DBName, defaultPGPool); err != nil {
		return nil, nil, err
	}

	dbPool, err := postgres.ConnectPool(
		context.Background(),
		getPostgresConnString(dockerResource.GetPort("5432/tcp"), cfg.DBName),
		zerolog.Nop(),
	)
	if err != nil {
		return nil, nil, err
	}

	if err = migrations.RunMigrations(
		getPostgresConnString(dockerResource.GetPort("5432/tcp"), cfg.DBName),
	); err != nil {
		return nil, nil, fmt.Errorf("running migrations: %w", err)
	}

	// Creates a template database with no active connections to be used as template on new dynamic database creations.
	_, _ = dbPool.Exec(
		context.Background(),
		fmt.Sprintf("create database %[1]s_template TEMPLATE %[1]s;", dbPool.Config().ConnConfig.Database),
	)
	concurrentConn = dbPool

	teardownFn = func() {
		dbPool.Close()

		_ = dropDB(cfg.DBName, defaultPGPool)
		_ = dropDB(fmt.Sprintf("%s_template", cfg.DBName), defaultPGPool)

		defaultPGPool.Close()

		_ = dockerResource.Close()
	}

	return &DockerizedPostgres{
		Pool: dbPool,
		Port: dbPort,
	}, teardownFn, nil
}

func getDockerPostgresResource(dockerPool *dockertest.Pool, config DockerContainerConfig) (*dockertest.Resource, error) {
	var containerName string

	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Name:       containerName,
		Repository: "postgres",
		Tag:        config.Version,
		Env:        []string{"POSTGRES_USER=postgres", "POSTGRES_PASSWORD=postgres", "POSTGRES_DB=" + config.DBName},
	}, func(c *docker.HostConfig) {
		// Set AutoRemove to true so that stopped container goes away by itself.
		c.AutoRemove = true
		c.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, fmt.Errorf(`RunWithOptions failed: %w`, err)
	}

	return resource, nil
}

func pingPostgresFn(port string) func() error {
	return func() error {
		conn, err := postgres.Connect(
			context.Background(),
			getPostgresConnString(port, "postgres"),
			zerolog.Nop(),
		)
		if err != nil {
			return err
		}

		defer conn.Close(context.Background())
		return conn.Ping(context.Background())
	}
}

func dropDB(dbName string, pool *pgxpool.Pool) error {
	if dbName == _postgresDefaultDB {
		return nil // can't drop default db
	}

	if _, err := pool.Exec(context.Background(), fmt.Sprintf(`DROP DATABASE IF EXISTS %s;`, dbName)); err != nil {
		return fmt.Errorf("failed dropping database %s: %w", dbName, err)
	}

	return nil
}

func createDB(dbName string, pool *pgxpool.Pool) error {
	if dbName == _postgresDefaultDB {
		return nil // can't create default db
	}

	_ = dropDB(dbName, pool)

	if _, err := pool.Exec(context.Background(), fmt.Sprintf(`CREATE DATABASE %s;`, dbName)); err != nil {
		return fmt.Errorf("error creating database %s: %w", dbName, err)
	}

	return nil
}

func getPostgresConnString(port, dbName string) string {
	return fmt.Sprintf("postgres://postgres:postgres@localhost:%s/%s?sslmode=disable", port, dbName)
}
