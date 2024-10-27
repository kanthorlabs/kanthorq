package core

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xvalidator"
)

//go:embed stream_get_event.sql
var StreamGetEventSql string

type StreamGetEventReq struct {
	Stream   *entities.StreamRegistry `validate:"required"`
	EventIds []string                 `validate:"required,gt=0,lte=500,dive,required"`
}

type StreamGetEventRes struct {
	Events []*entities.Event
}

func (req *StreamGetEventReq) Do(ctx context.Context, tx pgx.Tx) (*StreamGetEventRes, error) {
	err := xvalidator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(req.EventIds))
	args := pgx.NamedArgs{}
	for i, id := range req.EventIds {
		binding := fmt.Sprintf("event_id_%d", i)
		names[i] = fmt.Sprintf("@%s", binding)
		args[binding] = id
	}
	table := pgx.Identifier{entities.Collection(req.Stream.Id)}.Sanitize()
	query := fmt.Sprintf(StreamGetEventSql, table, strings.Join(names, ","))
	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := &StreamGetEventRes{Events: make([]*entities.Event, 0)}
	for rows.Next() {
		var event entities.Event
		err := rows.Scan(
			&event.Id,
			&event.Subject,
			&event.Body,
			&event.Metadata,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		res.Events = append(res.Events, &event)
	}

	// rows.Err returns any error that occurred while reading
	// always check it before finishing the read
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
