package core

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/stretchr/testify/require"
)

func Seed(
	t *testing.T,
	ctx context.Context,
	conn *pgx.Conn,
) (*entities.StreamRegistry, *entities.ConsumerRegistry) {
	req := &ConsumerRegisterReq{
		StreamName:            xfaker.StreamName(),
		ConsumerName:          xfaker.ConsumerName(),
		ConsumerSubjectFilter: []string{xfaker.Subject()},
		ConsumerAttemptMax:    xfaker.F.Int16Between(2, 10),
	}
	// ConsumerRegister also register stream
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	return res.StreamRegistry, res.ConsumerRegistry
}

func SeedEvents(
	t *testing.T,
	ctx context.Context,
	conn *pgx.Conn,
	stream *entities.StreamRegistry,
	consumer *entities.ConsumerRegistry,
	count int,
) []*entities.Event {
	events := FakeEvents(xfaker.SubjectWihtPattern(consumer.SubjectFilter[0]), count)

	req := &StreamPutEventsReq{
		Stream: stream,
		Events: events,
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)
	require.Equal(t, int64(count), res.InsertCount)

	return events
}

func SeedTasks(
	t *testing.T,
	ctx context.Context,
	conn *pgx.Conn,
	consumer *entities.ConsumerRegistry,
	events []*entities.Event,
	state entities.TaskState,
) []*entities.Task {
	tasks := FakeTasks(events, state)

	req := &ConsumerPutTasksReq{
		Consumer: consumer,
		Tasks:    tasks,
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)
	require.Equal(t, int64(len(tasks)), res.InsertCount)

	return tasks
}

func FakeEvents(subject string, count int) []*entities.Event {
	events := make([]*entities.Event, count)
	for i := 0; i < count; i++ {
		events[i] = entities.NewEvent(subject, xfaker.DataOf16Kb())
	}
	return events
}

func FakeTasks(events []*entities.Event, state entities.TaskState) []*entities.Task {
	tasks := make([]*entities.Task, len(events))
	for i := range events {
		tasks[i] = &entities.Task{
			EventId:      events[i].Id,
			Subject:      events[i].Subject,
			State:        int16(state),
			ScheduleAt:   time.Now().UTC().UnixMilli(),
			AttemptCount: 1,
			AttemptedAt:  time.Now().UTC().UnixMilli(),
			FinalizedAt:  0,
			CreatedAt:    events[i].CreatedAt,
			UpdatedAt:    time.Now().UTC().UnixMilli(),
		}
	}
	return tasks
}
