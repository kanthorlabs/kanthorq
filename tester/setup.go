package tester

import (
	"context"
	"errors"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

func SetupPostgres(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("KANTHORQ_POSTGRES_URI"))
	if err != nil {
		return nil, err
	}

	m, err := migrate.New(os.Getenv("KANTHORQ_MIGRATION_SOURCE"), os.Getenv("KANTHORQ_POSTGRES_URI"))
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	defer m.Close()

	return conn, nil
}
