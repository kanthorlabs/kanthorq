package kanthorq

import "github.com/kanthorlabs/kanthorq/pkg/faker"

func FakeEvents(topic string, min, max int64) []*Event {
	count := faker.F.Int64Between(min, max)
	events := make([]*Event, count)
	for i := 0; i < int(count); i++ {
		events[i] = NewEvent(topic, faker.DataOf16Kb())
	}
	return events
}
