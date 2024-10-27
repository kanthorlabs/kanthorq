package core

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func Seed(
	t *testing.T,
	ctx context.Context,
	conn *pgx.Conn,
) (*entities.StreamRegistry, *entities.ConsumerRegistry) {
	req := &ConsumerRegisterReq{
		StreamName:                xfaker.StreamName(),
		ConsumerName:              xfaker.ConsumerName(),
		ConsumerSubjectIncludes:   []string{xfaker.Subject()},
		ConsumerSubjectExcludes:   []string{xfaker.Subject()},
		ConsumerAttemptMax:        xfaker.F.Int16Between(2, 10),
		ConsumerVisibilityTimeout: xfaker.F.Int64Between(15000, 300000),
	}
	// ConsumerRegister also register stream
	res, err := Do(ctx, conn, req)
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
	events := tester.FakeEvents(xfaker.SubjectWihtPattern(consumer.SubjectIncludes[0]), count)

	req := &StreamPutEventsReq{
		Stream: stream,
		Events: events,
	}
	res, err := Do(ctx, conn, req)
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
	tasks := tester.FakeTasks(events, state)

	req := &ConsumerPutTasksReq{
		Consumer: consumer,
		Tasks:    tasks,
	}
	res, err := Do(ctx, conn, req)
	require.NoError(t, err)
	require.Equal(t, int64(len(tasks)), res.InsertCount)

	return tasks
}
