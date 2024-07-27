package kanthorq

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/pkg/faker"
	"github.com/stretchr/testify/require"
)

func FakeEntities(t *testing.T, ctx context.Context, conn *pgx.Conn, count int) (*StreamRegistry, *ConsumerRegistry, []*Event) {
	creq := &ConsumerRegisterReq{
		StreamName:         faker.StreamName(),
		ConsumerName:       faker.ConsumerName(),
		ConsumerTopic:      faker.Topic(),
		ConsumerAttemptMax: faker.F.Int16Between(1, 10),
	}
	// ConsumerRegister also register stream
	cres, err := Do(ctx, creq, conn)
	require.NoError(t, err)

	events := FakeEvents(creq.ConsumerTopic, count)

	// put events to stream
	sreq := &StreamPutEventsReq{
		Stream: cres.StreamRegistry,
		Events: events,
	}
	sres, err := StreamPutEvents(ctx, sreq, conn)
	require.NoError(t, err)
	require.Equal(t, int64(count), sres.InsertCount)

	return cres.StreamRegistry, cres.ConsumerRegistry, events
}

func FakeEvents(topic string, count int) []*Event {
	events := make([]*Event, count)
	for i := 0; i < count; i++ {
		events[i] = NewEvent(topic, faker.DataOf16Kb())
	}
	return events
}
