package pub

import (
	"fmt"
	"strings"
	"time"

	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/pkg/command"
	"github.com/kanthorlabs/kanthorq/pkg/faker"
	"github.com/spf13/pflag"
)

func GetBody(flags *pflag.FlagSet) []byte {
	data, err := flags.GetString("body")
	if err != nil {
		panic(err)
	}

	if data == "__KANTHORQ_FAKE.DATA_OF_16KB__" {
		return faker.DataOf16Kb()
	}

	return []byte(data)
}

func GetMetadata(flags *pflag.FlagSet) kanthorq.Metadata {
	data, err := flags.GetStringArray("metadata")
	if err != nil {
		panic(err)
	}

	metadata := kanthorq.Metadata{}
	for i := 0; i < len(data); i++ {
		kv := strings.Split(data[i], "=")
		metadata[kv[0]] = kv[1]
	}

	return metadata
}

func GetEvents(flags *pflag.FlagSet) []*kanthorq.Event {
	subjectOrPattern := command.GetString(flags, "subject")
	body := GetBody(flags)
	metadata := GetMetadata(flags)

	count := command.GetInt(flags, "count")
	events := make([]*kanthorq.Event, count)
	for i := 0; i < count; i++ {
		subject := faker.SubjectWihtPattern(subjectOrPattern)
		event := kanthorq.NewEvent(subject, body)
		event.Metadata.Merge(metadata)
		event.Metadata["index"] = i
		events[i] = event

		ts := time.UnixMilli(event.CreatedAt).Format(time.RFC3339)
		fmt.Printf("%s | %s | %s\n", event.Id, event.Subject, ts)
	}

	return events
}
