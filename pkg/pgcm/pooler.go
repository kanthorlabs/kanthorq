package pgcm

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
)

// NewPooler initializes a ConnectionManager that connect to an external PG pooler
// PGBouncer for instance
// connections should be handled by the pooler instead of our client code
// so everytime we finish with a connection, we should return it to the pooler by closing it
func NewPooler(uri string) ConnectionManager {
	return &pooler{uri: uri}
}

type pooler struct {
	uri string

	mu sync.Mutex
}

func (cm *pooler) Start(ctx context.Context) error {
	return nil
}

func (cm *pooler) Stop(ctx context.Context) error {
	return nil
}

func (cm *pooler) Connection(ctx context.Context) (Connection, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conn, err := pgx.Connect(ctx, cm.uri)
	if err != nil {
		return nil, err
	}

	return &poolerc{conn: conn}, nil
}
