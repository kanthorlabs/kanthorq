package base

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed .version
var version string

func New() *cobra.Command {
	command := &cobra.Command{
		Short: short,
		Long:  long + short,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			}
		},
	}

	command.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(long + version + "\n")
		},
	})

	command.PersistentFlags().BoolP("verbose", "v", false, "show verbose output including debug information")

	return command
}
