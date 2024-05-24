package testify

import (
	"context"
	"errors"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/common/containers"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func SpinPostgres(ctx context.Context, name string) (*postgres.PostgresContainer, error) {
	// spin a container
	container, err := containers.Postgres(ctx, name)
	if err != nil {
		return nil, err
	}
	uri, err := containers.PostgresConnectionString(context.Background(), container)
	if err != nil {
		return nil, err
	}

	// run a migration
	m, err := migrate.New(os.Getenv("TEST_MIGRATION_SOURCE"), uri)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); errors.Is(err, migrate.ErrNoChange) || errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	defer m.Close()

	conn, err := pgxpool.New(context.Background(), uri)
	if err != nil {
		return nil, err
	}

	// cleanup
	if err := QueryTruncateConsumer()(ctx, conn); err != nil {
		return nil, err
	}
	if err := QueryTruncateStream()(ctx, conn); err != nil {
		return nil, err
	}

	return container, nil
}
