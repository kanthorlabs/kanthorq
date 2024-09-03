package puller

import (
	"context"
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq/core"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/pgcm"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestPrimary_Do(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	cm, err := pgcm.New(os.Getenv("KANTHORQ_POSTGRES_URI"))
	require.NoError(t, err)

	res, err := core.DoWithCM(ctx, &core.ConsumerRegisterReq{
		StreamName:                xfaker.StreamName(),
		ConsumerName:              xfaker.ConsumerName(),
		ConsumerSubjectFilter:     []string{xfaker.Subject()},
		ConsumerAttemptMax:        xfaker.F.Int16Between(2, 10),
		ConsumerVisibilityTimeout: xfaker.F.Int64Between(15000, 300000),
	}, cm)
	require.NoError(t, err)

	// need 2 batches of events to pull
	count := xfaker.F.IntBetween(101, 199)
	events := tester.FakeEvents(xfaker.SubjectWihtPattern(res.ConsumerRegistry.SubjectFilter[0]), count)

	_, err = core.DoWithCM(ctx, &core.StreamPutEventsReq{
		Stream: res.StreamRegistry,
		Events: events,
	}, cm)
	require.NoError(t, err)

	in := PullerIn{
		Size:        100,
		WaitingTime: 3000,
	}
	p := New(cm, res.StreamRegistry, res.ConsumerRegistry, in)

	first, err := p.Do(ctx)
	require.NoError(t, err)

	// first batch should return enough tasks
	require.Equal(t, in.Size, len(first.Tasks))
	require.Equal(t, in.Size, len(first.Events))
	require.Equal(t, in.Size, len(first.EventIds))

	// all tasks must be at running state
	for _, task := range first.Tasks {
		require.Equal(t, entities.StateRunning, task.State)
	}

	second, err := p.Do(ctx)
	require.NoError(t, err)

	// the second batch should return remaining tasks
	require.Equal(t, count-in.Size, len(second.Tasks))
	require.Equal(t, count-in.Size, len(second.Events))
	require.Equal(t, count-in.Size, len(second.EventIds))

	// all tasks must be at running state
	for _, task := range second.Tasks {
		require.Equal(t, entities.StateRunning, task.State)
	}

}

func TestPrimary_Do_NoEvent(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	defer func() {
		require.NoError(t, conn.Close(ctx))
	}()
	require.NoError(t, err)

	cm, err := pgcm.New(os.Getenv("KANTHORQ_POSTGRES_URI"))
	require.NoError(t, err)

	res, err := core.DoWithCM(ctx, &core.ConsumerRegisterReq{
		StreamName:                xfaker.StreamName(),
		ConsumerName:              xfaker.ConsumerName(),
		ConsumerSubjectFilter:     []string{xfaker.Subject()},
		ConsumerAttemptMax:        xfaker.F.Int16Between(2, 10),
		ConsumerVisibilityTimeout: xfaker.F.Int64Between(15000, 300000),
	}, cm)
	require.NoError(t, err)

	in := PullerIn{
		Size:        100,
		WaitingTime: 3000,
	}
	p := New(cm, res.StreamRegistry, res.ConsumerRegistry, in)

	out, err := p.Do(ctx)
	require.NoError(t, err)

	require.Empty(t, out.Tasks)
	require.Empty(t, out.Events)
	require.Empty(t, out.EventIds)
}
