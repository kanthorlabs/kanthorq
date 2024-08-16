package kanthorq

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/faker"
	"github.com/stretchr/testify/require"
)

func Seed(t *testing.T, ctx context.Context, conn *pgx.Conn) (*StreamRegistry, *ConsumerRegistry) {
	req := &ConsumerRegisterReq{
		StreamName:            faker.StreamName(),
		ConsumerName:          faker.ConsumerName(),
		ConsumerSubjectFilter: faker.Subject(),
		ConsumerAttemptMax:    faker.F.Int16Between(2, 10),
	}
	// ConsumerRegister also register stream
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	return res.StreamRegistry, res.ConsumerRegistry
}

func SeedEvents(t *testing.T, ctx context.Context, conn *pgx.Conn, stream *StreamRegistry, consumer *ConsumerRegistry, count int) []*Event {
	events := FakeEvents(faker.SubjectWihtPattern(consumer.SubjectFilter), count)

	req := &StreamPutEventsReq{
		Stream: stream,
		Events: events,
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)
	require.Equal(t, int64(count), res.InsertCount)

	return events
}

func SeedTasks(t *testing.T, ctx context.Context, conn *pgx.Conn, consumer *ConsumerRegistry, events []*Event, state TaskState) []*Task {
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

func FakeEvents(subject string, count int) []*Event {
	events := make([]*Event, count)
	for i := 0; i < count; i++ {
		events[i] = NewEvent(subject, faker.DataOf16Kb())
	}
	return events
}

func FakeTasks(events []*Event, state TaskState) []*Task {
	tasks := make([]*Task, len(events))
	for i := range events {
		tasks[i] = &Task{
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
