package api

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/kanthorq/entities"
)

func ConsumerJobMarkRetry(consumer *entities.Consumer, reports map[string]error) *ConsumerJobMarkRetryReq {
	return &ConsumerJobMarkRetryReq{
		Consumer: consumer,
		Reports:  reports,
	}
}

//go:embed consumer_job_mark_retry.sql
var ConsumerJobMarkRetrySQL string

type ConsumerJobMarkRetryReq struct {
	Consumer *entities.Consumer
	Reports  map[string]error
}

type ConsumerJobMarkRetryRes struct {
	Status map[string]bool
}

func (req *ConsumerJobMarkRetryReq) Do(ctx context.Context, tx pgx.Tx, clock clock.Clock) (*ConsumerJobMarkRetryRes, error) {
	res := &ConsumerJobMarkRetryRes{Status: make(map[string]bool)}

	var names = make([]string, len(req.Reports))
	var args = pgx.NamedArgs{
		"complete_state": entities.StateRetryable,
		"running_state":  entities.StateRunning,
		"attempt_at":     clock.Now().UnixMilli(),
	}
	var i = 0
	for id := range req.Reports {
		// we assume that events are not able to update firstly
		// later if we can update its state, we can set it back to true
		res.Status[id] = false

		eventIdBind := fmt.Sprintf("event_id_%d", i)
		names[i] = fmt.Sprintf("@%s", eventIdBind)
		args[eventIdBind] = id

		i++
	}

	table := pgx.Identifier{entities.CollectionConsumerJob(req.Consumer.Name)}.Sanitize()
	query := fmt.Sprintf(ConsumerJobMarkRetrySQL, table, strings.Join(names, ","))

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		res.Status[id] = true
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}
