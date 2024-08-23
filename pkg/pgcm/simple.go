package pgcm

import (
	"context"
	"errors"
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
	if !cm.cmu.TryLock() {
		return nil, errors.New("connection already in use")
	}

	if cm.conn != nil && !cm.conn.IsClosed() {
		return cm.conn, nil
	}

	conn, err := pgx.Connect(ctx, cm.uri)
	if err != nil {
		return nil, err
	}
	cm.conn = conn

	return cm.conn, nil
}

func (cm *simple) Release(ctx context.Context, conn *pgx.Conn) error {
	cm.cmu.Unlock()
	// don't close the connection
	return nil
}
