package core

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/validator"
)

//go:embed stream_scan.sql
var StreamScanSql string

// TODO: using cursor with safe window (-1000ms)
type StreamScanReq struct {
	Stream   *entities.StreamRegistry   `validate:"required"`
	Consumer *entities.ConsumerRegistry `validate:"required"`

	Size        int   `validate:"required,gt=0"`
	WaitingTime int64 `validate:"gte=1000"`
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

	waitctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(req.WaitingTime))
	defer cancel()

	res := &StreamScanRes{Cursor: req.Consumer.Cursor}
	for len(res.Ids) < req.Size {
		prev := res.Cursor
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-waitctx.Done():
			return res, nil
		default:
			if err := req.scan(ctx, tx, res); err != nil {
				return nil, err
			}

			// if cursor has not changed, that mean there no new rows, wait for a while
			if prev == res.Cursor {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-waitctx.Done():
					return res, nil
				case <-time.After(time.Millisecond * 300):
					log.Println("waiting for new events...")
				}
			}
		}
	}

	return res, nil
}

func (req *StreamScanReq) scan(ctx context.Context, tx pgx.Tx, res *StreamScanRes) error {
	table := pgx.Identifier{entities.Collection(req.Stream.Id)}.Sanitize()
	query := fmt.Sprintf(StreamScanSql, table)
	args := pgx.NamedArgs{
		"cursor": res.Cursor,
		"size":   req.Size,
	}

	rows, err := tx.Query(ctx, query, args)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	defer rows.Close()

	checked := make(map[string]bool)
	for rows.Next() && len(res.Ids) < req.Size {
		var id, subject string
		if err := rows.Scan(&id, &subject); err != nil {
			return err
		}

		// override cursor with newest id
		res.Cursor = id

		// a subject was already checked
		exist, seen := checked[subject]
		if exist {
			if seen {
				res.Ids = append(res.Ids, id)
			}
			continue
		}

		checked[subject] = req.match(subject)
		if checked[subject] {
			res.Ids = append(res.Ids, id)
		}
	}

	// rows.Err returns any error that occurred while reading
	// always check it before finishing the read
	return rows.Err()
}

func (req *StreamScanReq) match(subject string) bool {
	for i := 0; i < len(req.Consumer.SubjectFilter); i++ {
		if MatchSubject(req.Consumer.SubjectFilter[i], subject) {
			return true
		}
	}

	return false
}
