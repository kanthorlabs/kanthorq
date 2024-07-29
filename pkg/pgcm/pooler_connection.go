package pgcm

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// poolerc is a implementation of Connection
// because our connections are managed by the pooler
// it should be returned back to the pooler after we are done with our tasks
// so the Close method will close the connection
type poolerc struct {
	conn *pgx.Conn
}

func (c *poolerc) Raw() *pgx.Conn {
	return c.conn
}

func (c *poolerc) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}
