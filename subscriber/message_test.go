package subscriber

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/kanthorlabs/kanthorq/core"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestMessage_Ack(t *testing.T) {
	events := tester.FakeEvents(xfaker.Subject(), 1)
	tasks := tester.FakeTasks(events, entities.StateRunning)

	ctx := context.Background()

	cm, err := pgcm.New(os.Getenv("KANTHORQ_POSTGRES_URI"))
	require.NoError(t, err)

	options := &core.ConsumerRegisterReq{
		StreamName:                xfaker.StreamName(),
		ConsumerName:              xfaker.ConsumerName(),
		ConsumerSubjectIncludes:   []string{xfaker.Subject()},
		ConsumerAttemptMax:        xfaker.F.Int16Between(2, 10),
		ConsumerVisibilityTimeout: xfaker.F.Int64Between(15000, 300000),
	}
	options.ConsumerSubjectIncludes = append(options.ConsumerSubjectIncludes, events[0].Subject)

	res, err := core.DoWithCM(ctx, options, cm)
	require.NoError(t, err)

	msg := &Message{
		Event:    events[0],
		Task:     tasks[0],
		cm:       cm,
		consumer: res.ConsumerRegistry,
	}
	require.NoError(t, msg.Ack(ctx))
	// call it twice should be safe
	require.NoError(t, msg.Ack(ctx))
}

func TestMessage_Ack_Error(t *testing.T) {
	events := tester.FakeEvents(xfaker.Subject(), 1)
	tasks := tester.FakeTasks(events, entities.StateRunning)

	ctx := context.Background()

	cm, err := pgcm.New(os.Getenv("KANTHORQ_POSTGRES_URI"))
	require.NoError(t, err)

	options := &core.ConsumerRegisterReq{
		StreamName:                xfaker.StreamName(),
		ConsumerName:              xfaker.ConsumerName(),
		ConsumerSubjectIncludes:   []string{xfaker.Subject()},
		ConsumerAttemptMax:        xfaker.F.Int16Between(2, 10),
		ConsumerVisibilityTimeout: xfaker.F.Int64Between(15000, 300000),
	}
	options.ConsumerSubjectIncludes = append(options.ConsumerSubjectIncludes, events[0].Subject)

	res, err := core.DoWithCM(ctx, options, cm)
	require.NoError(t, err)

	msg := &Message{
		Event:    events[0],
		Task:     tasks[0],
		cm:       cm,
		consumer: res.ConsumerRegistry,
	}
	// override with non-exist consumer, so ack should fail
	msg.consumer.Id = entities.ConsumerId()

	require.ErrorContains(t, msg.Ack(ctx), "SQLSTATE 42P01")
	// call it twice should be safe
	require.NoError(t, msg.Ack(ctx))
	// only allow call either ack or nack once
	require.ErrorIs(t, msg.Nack(ctx, errors.New(time.Now().Format(time.RFC3339Nano))), ErrInvalidMessageState)
}

func TestMessage_AckTx(t *testing.T) {
	events := tester.FakeEvents(xfaker.Subject(), 1)
	tasks := tester.FakeTasks(events, entities.StateRunning)

	ctx := context.Background()

	cm, err := pgcm.New(os.Getenv("KANTHORQ_POSTGRES_URI"))
	require.NoError(t, err)

	options := &core.ConsumerRegisterReq{
		StreamName:                xfaker.StreamName(),
		ConsumerName:              xfaker.ConsumerName(),
		ConsumerSubjectIncludes:   []string{xfaker.Subject()},
		ConsumerAttemptMax:        xfaker.F.Int16Between(2, 10),
		ConsumerVisibilityTimeout: xfaker.F.Int64Between(15000, 300000),
	}
	options.ConsumerSubjectIncludes = append(options.ConsumerSubjectIncludes, events[0].Subject)

	res, err := core.DoWithCM(ctx, options, cm)
	require.NoError(t, err)

	conn, err := cm.Acquire(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cm.Release(ctx, conn))
	}()

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)

	msg := &Message{
		Event:    events[0],
		Task:     tasks[0],
		cm:       cm,
		consumer: res.ConsumerRegistry,
	}
	require.NoError(t, msg.AckTx(ctx, tx))
	// call it twice should be safe
	require.NoError(t, msg.AckTx(ctx, tx))

	require.NoError(t, tx.Commit(ctx))
}

func TestMessage_Nack(t *testing.T) {
	events := tester.FakeEvents(xfaker.Subject(), 1)
	tasks := tester.FakeTasks(events, entities.StateRunning)

	ctx := context.Background()

	cm, err := pgcm.New(os.Getenv("KANTHORQ_POSTGRES_URI"))
	require.NoError(t, err)

	options := &core.ConsumerRegisterReq{
		StreamName:                xfaker.StreamName(),
		ConsumerName:              xfaker.ConsumerName(),
		ConsumerSubjectIncludes:   []string{xfaker.Subject()},
		ConsumerAttemptMax:        xfaker.F.Int16Between(2, 10),
		ConsumerVisibilityTimeout: xfaker.F.Int64Between(15000, 300000),
	}
	options.ConsumerSubjectIncludes = append(options.ConsumerSubjectIncludes, events[0].Subject)

	res, err := core.DoWithCM(ctx, options, cm)
	require.NoError(t, err)

	msg := &Message{
		Event:    events[0],
		Task:     tasks[0],
		cm:       cm,
		consumer: res.ConsumerRegistry,
	}

	fakerr := errors.New(time.Now().Format(time.RFC3339Nano))
	require.NoError(t, msg.Nack(ctx, fakerr))
	// call it twice should be safe
	require.NoError(t, msg.Nack(ctx, fakerr))
}

func TestMessage_Nack_Error(t *testing.T) {
	events := tester.FakeEvents(xfaker.Subject(), 1)
	tasks := tester.FakeTasks(events, entities.StateRunning)

	ctx := context.Background()

	cm, err := pgcm.New(os.Getenv("KANTHORQ_POSTGRES_URI"))
	require.NoError(t, err)

	options := &core.ConsumerRegisterReq{
		StreamName:                xfaker.StreamName(),
		ConsumerName:              xfaker.ConsumerName(),
		ConsumerSubjectIncludes:   []string{xfaker.Subject()},
		ConsumerAttemptMax:        xfaker.F.Int16Between(2, 10),
		ConsumerVisibilityTimeout: xfaker.F.Int64Between(15000, 300000),
	}
	options.ConsumerSubjectIncludes = append(options.ConsumerSubjectIncludes, events[0].Subject)

	res, err := core.DoWithCM(ctx, options, cm)
	require.NoError(t, err)

	msg := &Message{
		Event:    events[0],
		Task:     tasks[0],
		cm:       cm,
		consumer: res.ConsumerRegistry,
	}
	// override with non-exist consumer, so ack should fail
	msg.consumer.Id = entities.ConsumerId()

	ferr := errors.New(time.Now().Format(time.RFC3339Nano))

	require.ErrorContains(t, msg.Nack(ctx, ferr), "SQLSTATE 42P01")
	// call it twice should be safe
	require.NoError(t, msg.Nack(ctx, ferr))
	// only allow call either ack or nack once
	require.ErrorIs(t, msg.Ack(ctx), ErrInvalidMessageState)
}

func TestMessage_NackTx(t *testing.T) {
	events := tester.FakeEvents(xfaker.Subject(), 1)
	tasks := tester.FakeTasks(events, entities.StateRunning)

	ctx := context.Background()

	cm, err := pgcm.New(os.Getenv("KANTHORQ_POSTGRES_URI"))
	require.NoError(t, err)

	options := &core.ConsumerRegisterReq{
		StreamName:                xfaker.StreamName(),
		ConsumerName:              xfaker.ConsumerName(),
		ConsumerSubjectIncludes:   []string{xfaker.Subject()},
		ConsumerAttemptMax:        xfaker.F.Int16Between(2, 10),
		ConsumerVisibilityTimeout: xfaker.F.Int64Between(15000, 300000),
	}
	options.ConsumerSubjectIncludes = append(options.ConsumerSubjectIncludes, events[0].Subject)

	res, err := core.DoWithCM(ctx, options, cm)
	require.NoError(t, err)

	conn, err := cm.Acquire(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cm.Release(ctx, conn))
	}()

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)

	msg := &Message{
		Event:    events[0],
		Task:     tasks[0],
		cm:       cm,
		consumer: res.ConsumerRegistry,
	}

	fakerr := errors.New(time.Now().Format(time.RFC3339Nano))
	require.NoError(t, msg.NackTx(ctx, fakerr, tx))
	// call it twice should be safe
	require.NoError(t, msg.NackTx(ctx, fakerr, tx))

	require.NoError(t, tx.Commit(ctx))
}
