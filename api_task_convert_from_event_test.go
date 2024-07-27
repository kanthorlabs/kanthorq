package kanthorq

import (
	"context"
	"testing"

	"github.com/kanthorlabs/kanthorq/pkg/faker"
	"github.com/kanthorlabs/kanthorq/tester"
	"github.com/stretchr/testify/require"
)

func TestTaskConvertFromEvent(t *testing.T) {
	ctx := context.Background()
	conn, err := tester.SetupPostgres(ctx)
	require.NoError(t, err)

	topic := faker.Topic()

	// ConsumerRegister also register stream
	registry, err := Do(ctx, &ConsumerRegisterReq{
		StreamName:         faker.StreamName(),
		ConsumerName:       faker.ConsumerName(),
		ConsumerTopic:      topic,
		ConsumerAttemptMax: faker.F.Int16Between(1, 10),
	}, conn)
	require.NoError(t, err)

	// insert events
	events := FakeEvents(topic, 100, 500)
	size := len(events) - 1

	_, err = StreamPutEvents(ctx, &StreamPutEventsReq{
		Stream: registry.StreamRegistry,
		Events: events,
	}, conn)
	require.NoError(t, err)

	req := &TaskConvertFromEventReq{
		ConsumerName:     registry.ConsumerRegistry.Name,
		Size:             size,
		InitialTaskState: StateAvailable,
	}
	res, err := Do(ctx, req, conn)
	require.NoError(t, err)

	require.Equal(t, size, len(res.Tasks))
}
