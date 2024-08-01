package pub

import (
	"strings"

	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/pkg/faker"
	"github.com/spf13/pflag"
)

func GetBody(flags *pflag.FlagSet) []byte {
	data, err := flags.GetString("body")
	if err != nil {
		panic(err)
	}

	if data == "kanthorq.fake.16kb" {
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
