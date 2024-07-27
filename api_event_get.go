package kanthorq

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

//go:embed api_event_get.sql
var EventGetSql string

type EventGetReq struct {
	Stream   *StreamRegistry `validate:"required"`
	EventIds []string        `validate:"required,gt=0,dive,required"`
}

type EventGetRes struct {
	Events []*Event
}

func (req *EventGetReq) Do(ctx context.Context, tx pgx.Tx) (*EventGetRes, error) {
	err := validator.Validate.Struct(req)
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
	table := pgx.Identifier{Collection(req.Stream.Id)}.Sanitize()
	query := fmt.Sprintf(EventGetSql, table, strings.Join(names, ","))
	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := &EventGetRes{Events: make([]*Event, 0)}
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.Id,
			&event.Topic,
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
