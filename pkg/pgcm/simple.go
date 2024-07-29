package pgcm

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
)

// NewSimple initializes a simple implementation of ConnectionManager
// it holds a single connection to the database at once
// if the connection is closed, it will create a new one
func NewSimple(uri string) ConnectionManager {
	return &simple{uri: uri}
}

type simple struct {
	uri string

	mu   sync.Mutex
	conn *pgx.Conn
}

func (cm *simple) Start(ctx context.Context) error {
	return nil
}

func (cm *simple) Stop(ctx context.Context) error {
	if cm.conn == nil {
		return nil
	}

	if cm.conn.IsClosed() {
		return nil
	}

	return cm.conn.Close(ctx)
}

func (cm *simple) Connection(ctx context.Context) (Connection, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.conn != nil && !cm.conn.IsClosed() {
		return &simplec{conn: cm.conn}, nil
	}

	conn, err := pgx.Connect(ctx, cm.uri)
	if err != nil {
		return nil, err
	}
	cm.conn = conn

	return &simplec{conn: cm.conn}, nil
}
