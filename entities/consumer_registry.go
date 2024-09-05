package entities

import (
	"time"

	"github.com/kanthorlabs/kanthorq/pkg/xid"
)

func ConsumerId() string {
	return xid.New("consumer")
}

func ConsumerIdFromTime(t time.Time) string {
	return xid.NewWithTime("consumer", t)
}

type ConsumerRegistry struct {
	StreamId          string   `json:"stream_id" validate:"required"`
	StreamName        string   `json:"stream_name" validate:"required,is_collection_name"`
	Id                string   `json:"id" validate:"required"`
	Name              string   `json:"name" validate:"required,is_collection_name"`
	SubjectIncludes   []string `json:"subject_includes" validate:"required,gt=0,dive,is_subject_filter"`
	SubjectExcludes   []string `json:"subject_excludes" validate:"gte=0,dive,is_subject_filter"`
	Cursor            string   `json:"cursor"`
	AttemptMax        int16    `json:"attempt_max"`
	VisibilityTimeout int64    `json:"visibility_timeout" validate:"required,gt=1000"`
	CreatedAt         int64    `json:"created_at"`
	UpdatedAt         int64    `json:"updated_at"`
}
