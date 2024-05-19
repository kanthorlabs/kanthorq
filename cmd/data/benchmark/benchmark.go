package benchmark

import (
	"github.com/spf13/cobra"
)

var namespace = "benchmark"

func New() *cobra.Command {
	command := &cobra.Command{
		Use: "benchmark",
	}

	command.AddCommand(NewCleanup())
	command.AddCommand(NewPrepare())
	command.AddCommand(NewSeed())

	command.PersistentFlags().StringP("storage", "", "", "path to your benchmark data storage")
	command.MarkPersistentFlagRequired("storage")
	command.PersistentFlags().StringP("connection", "", "", "database connections string")
	command.MarkPersistentFlagRequired("connection")
	return command
}
