package pgcm

import (
	"context"
	"net/url"

	"github.com/jackc/pgx/v5"
)

// New returns a new ConnectionManager based on the connection string
// if the connection contains the default_query_exec_mode query parameter
// we implicitly assume that the client want to connect to a pooler
func New(uri string) (ConnectionManager, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	isPooler := u.Query().Has("pooling")

	// must remove our custom parameter, otherwise it will be rejected
	query := u.Query()
	query.Del("pooling")
	u.RawQuery = query.Encode()

	if isPooler {
		return NewPooler(u.String()), nil
	}

	return NewSimple(u.String()), nil
}

type ConnectionManager interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Connection(ctx context.Context) (Connection, error)
}

type Connection interface {
	Raw() *pgx.Conn
	Close(ctx context.Context) error
}
