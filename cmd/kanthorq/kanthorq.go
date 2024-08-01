package main

import (
	_ "embed"

	"github.com/kanthorlabs/kanthorq/cmd/kanthorq/migrate"
	"github.com/kanthorlabs/kanthorq/cmd/kanthorq/pub"
	"github.com/kanthorlabs/kanthorq/cmd/kanthorq/show"
	"github.com/spf13/cobra"
)

var (
	//go:embed .version
	version string
	//go:embed .text-to-ascii
	logo    string
	tagline string = "Message Broker backed by PostgreSQL"
)

func New() *cobra.Command {
	command := &cobra.Command{
		Short: tagline,
		Long:  tagline + "\n" + logo,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			}
		},
	}

	command.AddCommand(show.New(logo, version))
	command.AddCommand(migrate.New())
	command.AddCommand(pub.New())

	command.PersistentFlags().BoolP("verbose", "v", false, "show verbose output including debug information")

	return command
}
