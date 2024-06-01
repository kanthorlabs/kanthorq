package kanthorq

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// IMPORTANT NOTES
//   - ALWAYS set pool_min_conns and pool_max_conns. pool_min_conns should be to
//     reserve one connection to prevent cold start problem
//   - ALWAYS set pool_max_conn_lifetime to a short duration to avoid the client
//     keep a connection is too long. The longer lifetime you set, the larger memory of
//     that connection consumes
//   - SHOULD set pool_max_conn_idle_time to release free connections
//
// DEFAULT CONFIG
//   - pool_min_conns = 0
//   - pool_max_conns = 4
//   - pool_max_conn_lifetime = time.Hour
//   - pool_max_conn_idle_time = time.Minute * 30
func Connection(ctx context.Context, uri string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, uri)
}
