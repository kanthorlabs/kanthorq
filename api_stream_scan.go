package kanthorq

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

//go:embed api_stream_scan.sql
var StreamScanSql string

type StreamScanReq struct {
	Stream   *StreamRegistry   `validate:"required"`
	Consumer *ConsumerRegistry `validate:"required"`

	Size        int `validate:"required,gt=0"`
	IntervalMax int `validate:"gt=0"`
}

type StreamScanRes struct {
	Ids    []string
	Cursor string
}

func (req *StreamScanReq) Do(ctx context.Context, tx pgx.Tx) (*StreamScanRes, error) {
	err := validator.Validate.Struct(req)
	if err != nil {
		return nil, err
	}

	var res = &StreamScanRes{Cursor: req.Consumer.Cursor}
	var interval = 0
	for len(res.Ids) < req.Size && interval < req.IntervalMax {
		if err := req.scan(ctx, tx, res); err != nil {
			return nil, err
		}
		interval++
	}

	return res, nil
}

func (req *StreamScanReq) scan(ctx context.Context, tx pgx.Tx, res *StreamScanRes) error {
	table := pgx.Identifier{Collection(req.Stream.Id)}.Sanitize()
	query := fmt.Sprintf(StreamScanSql, table)
	args := pgx.NamedArgs{
		"cursor": res.Cursor,
		"size":   req.Size,
	}

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() && len(res.Ids) < req.Size {
		var id, subject string
		if err := rows.Scan(&id, &subject); err != nil {
			return err
		}

		if MatchSubject(req.Consumer.Subject, subject) {
			res.Ids = append(res.Ids, id)
		}

		// override cursor with newest id
		res.Cursor = id
	}

	// rows.Err returns any error that occurred while reading
	// always check it before finishing the read
	return rows.Err()
}
