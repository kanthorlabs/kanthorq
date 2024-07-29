package pgcm

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// simplec is a implementation of Connection
// it does not nothing when you call the Close method
// because the connection is sharing across multiple places
// it should be closed only by the manager
type simplec struct {
	conn *pgx.Conn
}

func (c *simplec) Raw() *pgx.Conn {
	return c.conn
}

func (c *simplec) Close(ctx context.Context) error {
	return nil
}
