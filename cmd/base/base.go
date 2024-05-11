package base

import (
	_ "embed"

	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/project"
	"github.com/spf13/cobra"
)

//go:embed .version
var version string

func New() (configuration.Provider, *cobra.Command) {
	project.SetVersion(version)

	provider, err := configuration.New(project.Namespace())
	if err != nil {
		panic(err)
	}

	command := &cobra.Command{
		Short: short,
		Long:  long,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			}
		},
	}

	command.PersistentFlags().BoolP("verbose", "v", false, "show verbose output including debug information")

	return provider, command
}
