package pub

import (
	"fmt"
	"strings"
	"time"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xcmd"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/spf13/pflag"
)

func GetBody(flags *pflag.FlagSet) []byte {
	data, err := flags.GetString("body")
	if err != nil {
		panic(err)
	}

	if data == "__KANTHORQ_FAKE__.__DATA_OF_16KB__" {
		return xfaker.DataOf16Kb()
	}

	return []byte(data)
}

func GetMetadata(flags *pflag.FlagSet) entities.Metadata {
	data, err := flags.GetStringArray("metadata")
	if err != nil {
		panic(err)
	}

	metadata := entities.Metadata{}
	for i := 0; i < len(data); i++ {
		kv := strings.Split(data[i], "=")
		metadata[kv[0]] = kv[1]
	}

	return metadata
}

func GetEvents(flags *pflag.FlagSet) []*entities.Event {
	subjectOrPattern := xcmd.GetString(flags, "subject")
	body := GetBody(flags)
	metadata := GetMetadata(flags)

	count := xcmd.GetInt(flags, "count")
	events := make([]*entities.Event, count)
	for i := 0; i < count; i++ {
		subject := xfaker.SubjectWihtPattern(subjectOrPattern)
		event := entities.NewEvent(subject, body)
		event.Metadata.Merge(metadata)
		event.Metadata["index"] = i
		events[i] = event

		ts := time.UnixMilli(event.CreatedAt).Format(time.RFC3339)
		fmt.Printf("%s | %s | %s\n", event.Id, event.Subject, ts)
	}

	return events
}
