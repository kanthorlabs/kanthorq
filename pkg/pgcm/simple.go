package pgcm

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
)

var _ ConnectionManager = (*simple)(nil)

// NewSimple initializes a simple implementation of ConnectionManager
// it holds a single connection to the database at once
// if the connection is closed, it will create a new one
func NewSimple(uri string) ConnectionManager {
	return &simple{uri: uri}
}

type simple struct {
	uri string

	mu   sync.Mutex
	cmu  sync.Mutex
	conn *pgx.Conn
}

func (cm *simple) Start(ctx context.Context) error {
	return nil
}

func (cm *simple) Stop(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.conn == nil {
		return nil
	}

	if cm.conn.IsClosed() {
		return nil
	}

	return cm.conn.Close(ctx)
}

func (cm *simple) Acquire(ctx context.Context) (*pgx.Conn, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	cm.cmu.Lock()

	if cm.conn != nil && !cm.conn.IsClosed() {
		return cm.conn, nil
	}

	conn, err := pgx.Connect(ctx, cm.uri)
	if err != nil {
		// once we got an error, we should release the lock
		cm.cmu.Unlock()
		return nil, err
	}
	cm.conn = conn

	return cm.conn, nil
}

func (cm *simple) Release(ctx context.Context, conn *pgx.Conn) error {
	cm.cmu.Unlock()
	// don't close the connection
	return ctx.Err()
}
