package testify

import (
	"context"
	"errors"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, os.Getenv("TEST_DATABASE_URI"))
	if err != nil {
		return nil, err
	}

	// run a migration
	m, err := migrate.New(os.Getenv("TEST_MIGRATION_SOURCE"), os.Getenv("TEST_DATABASE_URI"))
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	defer m.Close()

	// cleanup
	if err := QueryTruncateConsumer()(ctx, conn); err != nil {
		return nil, err
	}
	if err := QueryTruncateStream()(ctx, conn); err != nil {
		return nil, err
	}

	return conn, nil
}
