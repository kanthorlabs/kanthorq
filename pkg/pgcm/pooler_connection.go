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

func (c *poolerc) Close(ctx context.Context) {
	// our main logic was executed successfully
	// so closing error should not cause of revert or rollback
	// log it here is enough
	if err := c.conn.Close(ctx); err != nil {
		// @TODO: log the error here
	}
}
