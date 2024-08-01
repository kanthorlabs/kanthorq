package command

import (
	"github.com/spf13/pflag"
)

func GetString(flags *pflag.FlagSet, name string) string {
	data, err := flags.GetString(name)
	if err != nil {
		panic(err)
	}

	return data
}

func GetInt(flags *pflag.FlagSet, name string) int {
	data, err := flags.GetInt(name)
	if err != nil {
		panic(err)
	}

	return data
}
