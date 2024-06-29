package q

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
)

func Stream(ctx context.Context, conn *pgx.Conn, stream *entities.Stream) (*entities.Stream, error) {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	r, err := NewStreamEnsure(stream.Name).Do(ctx, tx)
	if err != nil {
		return nil, err
	}

	return r.Stream, tx.Commit(ctx)
}
