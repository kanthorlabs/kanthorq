package entities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kanthorlabs/kanthorq/pkg/xid"
)

func TaskId() string {
	return xid.New("task")
}

func TaskIdFromTime(t time.Time) string {
	return xid.NewWithTime("task", t)
}

type Task struct {
	EventId string `json:"event_id" validate:"required"`
	Subject string `json:"subject" validate:"required,is_subject"`

	// State is the state of task like `available` or `completed`.
	State TaskState `json:"state" validate:"required,is_enum"`
	// ScheduledAt is when the task is scheduled to become available to be
	// worked. Tasks default to running immediately, but may be scheduled
	// for the future when they're inserted. They may also be scheduled for
	// later because they were snoozed or because they errored and have
	// additional retry attempts remaining.
	ScheduleAt int64 `json:"schedule_at" validate:"required,gt=0"`

	// AttemptCount is the attempt number of the task. Tasks are inserted at 0, the
	// number is incremented to 1 the first time work its worked, and may
	// increment further if it's either snoozed or errors.
	AttemptCount int16 `json:"attempt_count" validate:"gte=0"`
	// AttemptedAt is the time that the task was last worked. Starts out as NOW()
	// on a new insert.
	AttemptedAt int64 `json:"attempted_at" validate:"gte=0"`
	// AttemptedBy is the set of client IDs that have worked this task.
	AttemptedError []AttemptedError `json:"attempted_error"`
	// FinalizedAt is the time at which the task was "finalized", meaning it was
	// either completed successfully or errored for the last time such that
	// it'll no longer be retried.
	FinalizedAt int64 `json:"finalized_at" validate:"gte=0"`

	// Metadata carry additional information about the task.
	Metadata Metadata `json:"metadata" validate:"required"`
	// CreatedAt is when the task record was created.
	CreatedAt int64 `json:"created_at" validate:"required,gt=0"`
	// CreatedAt is when the task record was updated.
	UpdatedAt int64 `json:"updated_at" validate:"gte=0"`
}

type TaskState int16

const (
	// StateDiscarded is the state for tasks that have errored enough times
	// that they're no longer eligible to be retried. Manual user invention
	// is required for them to be tried again.
	StateDiscarded TaskState = -102
	// StateCancelled is the state for tasks that have been manually cancelled
	// by user request.
	StateCancelled TaskState = -101
	// StatePending is a state for tasks to be parked while waiting for some
	// external action before they can be worked. Tasks in pending will never be
	// worked or deleted unless moved out of this state by the user.
	StatePending TaskState = 0
	// StateAvailable is the state for tasks that are immediately eligible to
	// be worked.
	StateAvailable TaskState = 1
	// StateRunning is the state for tasks tasks which are actively running.
	StateRunning TaskState = 2
	// Completed is the state for tasks that have successfully run to
	// completion.
	StateCompleted TaskState = 101
	// StateRetryable is the state for tasks that have errored, but will be
	// retried.
	StateRetryable TaskState = 102
)

func (state TaskState) String() string {
	switch state {
	case StateDiscarded:
		return "discarded"
	case StateCancelled:
		return "cancelled"
	case StateAvailable:
		return "available"
	case StateRunning:
		return "running"
	case StateCompleted:
		return "completed"
	case StateRetryable:
		return "retryable"
	default:
		return ""
	}
}

type AttemptedError struct {
	// At is the time at which the error occurred.
	At int64 `json:"at" validate:"required,gt=0"`
	// Error contains the stringified error of an error returned from a job or a
	// panic value in case of a panic.
	Error string `json:"error" validate:"required"`
	// Stack contains a stack trace from a job that panicked. The trace is
	// produced by invoking `debug.Stack()`.
	Stack string `json:"stack" validate:"required"`
}

// Scan implements the sql.Scanner interface to scan a value from the database into the Metadata struct
func (err *AttemptedError) Scan(value interface{}) error {
	if data, ok := value.([]byte); ok {
		return json.Unmarshal(data, err)
	}

	if data, ok := value.(string); ok {
		return json.Unmarshal([]byte(data), err)
	}

	return fmt.Errorf("KANTHORQ.ATTEMPT_ERROR.SCAN.ERR: only string or []byte supported, got %T", value)
}

// Value implements the driver.Valuer interface to convert the Metadata struct to a value that can be stored in the database
func (err AttemptedError) Value() (driver.Value, error) {
	return json.Marshal(err)
}
